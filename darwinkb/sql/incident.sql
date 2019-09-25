-- ======================================================================
-- incident.sql handles the static & real time incidents feeds
-- ======================================================================
create schema if not exists darwin;

-- ======================================================================
-- An active incident
-- ======================================================================
create table if not exists darwin.incident
(
    id             serial                   not null,               -- internal id
    incidentNumber char(32)                 not null,               -- unique id at NRE, 32 hex ID
    participantRef text,                                            -- unique id of system issuing entry identifier
    created        timestamp with time zone not null,               -- when this incident was created at NRE
    planned        bool                     not null default false, -- was this planned
    cleared        bool                     not null default false, -- has this incident been cleared
    summary        text                     not null,               -- summary text
    validFrom      timestamp with time zone not null,               -- When this incident is valid from
    validTo        timestamp with time zone,                        -- Incident valid until this, forever if null
    version        bigint                   not null default 0,     -- unique version to this incident. Will always increase for each change
    primary key (id)
);

create unique index if not exists incident_num on darwin.incident (incidentNumber);
create index if not exists incident_ic on darwin.incident (id, cleared);

-- indices containing validity dates
create index if not exists incident_icf on darwin.incident (id, cleared, validFrom);
create index if not exists incident_icft on darwin.incident (id, cleared, validFrom, validTo);

-- ======================================================================
-- Table to hold the latest xml for this incident
-- ======================================================================
create table if not exists darwin.incident_xml
(
    id      bigint    not null references darwin.incident (id), -- FK to incident table
    xml     xml       not null,                                 -- Original xml
    updated timestamp not null,                                 -- When this entry was updated
    primary key (id)
);

-- ======================================================================
-- Table of Tocs
-- ======================================================================
create table if not exists darwin.incident_affects_toc
(
    ref  char(2) not null, -- Atoc code for the toc
    name name    not null, -- Public name of Operator
    primary key (ref)
);

-- ======================================================================
-- Table of toc's affected by an incident
-- ======================================================================
create table if not exists darwin.incident_affects
(
    id  bigint  not null references darwin.incident (id),              -- FK to incident table
    toc char(2) not null references darwin.incident_affects_toc (ref), -- FK to incident toc table
    primary key (id, toc)
);

create index if not exists incident_affects_i on darwin.incident_affects (id);
create index if not exists incident_affects_t on darwin.incident_affects (toc);

-- ======================================================================
-- updateIncident updates an incident, creating one as necessary
-- ======================================================================

create or replace function darwin.updateIncident(pxml xml)
    returns void
as
$$
declare
    ns       text[][] := ARRAY [
        ['ic', 'http://nationalrail.co.uk/xml/incident'],
        ['co', 'http://nationalrail.co.uk/xml/common']
        ];
    iid      bigint; -- internal ID for incident
    inum     char(32); -- incident number
    iver     bigint;
    axml     xml;
    bxml     xml;
    aref     char(2); -- affected toc ref
    aplanned bool; -- planned incident
    acleared bool; -- cleared incident
    rec      record;
begin
    -- Look for the incident, could be one of 2 root elements if from static or live feed
    axml = (xpath('//ic:PtIncident | //uk.co.nationalrail.xml.incident.PtIncidentStructure', pxml, ns))[1];

    inum = (xpath('//ic:IncidentNumber/text()', axml, ns))[1]::char(32);
    iver = (xpath('//ic:Version/text()', axml, ns))[1]::text::bigint;

    -- Check existence & if it's newer
    select into rec * from darwin.incident i where i.incidentNumber = inum;
    if found and rec.version <= iver then
        raise notice 'ignoring incident % version %', inum, iver;
        return;
    end if;

    raise notice 'got %', (xpath('//ic:ParticipantRef/text()', axml, ns))[1]::text;

    aplanned = (xpath('//ic:Planned/text()', axml, ns))[1]::text::bool;
    acleared = (xpath('//ic:Cleared/text()', axml, ns))[1]::text::bool;

    -- The main entry
    insert into darwin.incident (incidentnumber, participantref, created, summary, validfrom, validto, version, planned,
                                 cleared)
    values (inum,
            (xpath('//ic:ParticipantRef/text()', axml, ns))[1]::text,
            (xpath('//ic:CreationTime/text()', axml, ns))[1]::text::timestamp with time zone,
            (xpath('//ic:Summary/text()', axml, ns))[1]::text,
            (xpath('//ic:ValidityPeriod/co:StartTime/text()', axml, ns))[1]::text::timestamp with time zone,
            (xpath('//ic:ValidityPeriod/co:EndTime/text()', axml, ns))[1]::text::timestamp with time zone,
            iver,
            case aplanned when true then true else false end,
            case acleared when true then true else false end)
    on conflict (incidentNumber)
        do update set participantRef=excluded.participantref,
                      created=excluded.created,
                      summary=excluded.summary,
                      validFrom=excluded.validfrom,
                      validTo=excluded.validto,
                      version=excluded.version
    returning id into iid;

    -- the xml
    insert into darwin.incident_xml (id, xml, updated)
    values (iid, pxml, now())
    on conflict (id)
        do update set xml=excluded.xml,
                      updated=excluded.updated;

    -- Replace all affected operators with the new one
    delete from darwin.incident_affects where id = iid;
    foreach bxml in array xpath('//ic:Affects/ic:Operators/ic:AffectedOperator', axml, ns)
        loop
            aref = (xpath('//ic:OperatorRef/text()', bxml, ns))[1]::text;
            insert into darwin.incident_affects_toc (ref, name)
            values (aref,
                    (xpath('//ic:OperatorName/text()', bxml, ns))[1]::text)
            on conflict do nothing;

            insert into darwin.incident_affects (id, toc)
            values (iid, aref)
            on conflict do nothing;
        end loop;
end ;
$$ language plpgsql;

-- ======================================================================
-- updateIncidents updates all incidents from the static feed
-- ======================================================================

create or replace function darwin.updateIncidents(pxml xml)
    returns void
as
$$
declare
    ns   text[][] := ARRAY [
        ['ic', 'http://nationalrail.co.uk/xml/incident'],
        ['co', 'http://nationalrail.co.uk/xml/common']
        ];
    axml xml;
begin
    foreach axml in array (xpath('//ic:Incidents/ic:PtIncident', pxml, ns))
        loop
            perform darwin.updateIncident(axml);
        end loop;
end ;
$$ language plpgsql;

-- ======================================================================
-- getIncident returns details about a specific incident
-- ======================================================================

-- drop function darwin.getIncident(piid char(32));
create or replace function darwin.getIncident(piid char(32))
    returns jsonb
as
$$
declare
    axml xml;
begin
    select into axml x.xml
        from darwin.incident_xml x
    inner join darwin.incident i on x.id = i.id
    where i.incidentnumber = piid;
    if not found then
        return null;
    end if;

    return xml_to_json( axml );
end ;
$$ language plpgsql;

-- ======================================================================
-- getIncidents returns incidents for a toc or all if toc is null
-- ======================================================================

-- drop function darwin.getIncident(piid char(32));
create or replace function darwin.getIncidents(ptoc char(2))
    returns json
as
$$
select json_agg(row)
from (with incidents as (
    select distinct on (i.incidentNumber) i.incidentNumber,
                                          i.created,
                                          i.planned,
                                          i.summary,
                                          i.validfrom,
                                          i.validto,
                                          i.version
    from darwin.incident i
             inner join darwin.incident_affects a on i.id = a.id
             inner join darwin.incident_affects_toc t on a.toc = ptoc or ptoc is null
    where not i.cleared
      --and i.validfrom <= now()
      --and (i.validto is null or i.validto >= now())
    order by i.incidentnumber
)
      select *
      from incidents
      order by created desc
     ) row
$$ language sql;
