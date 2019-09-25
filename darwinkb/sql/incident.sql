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
    participantRef text                     not null,               -- unique id of system issuing entry identifier
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
    ns   text[][] := ARRAY [
        ['ic', 'http://nationalrail.co.uk/xml/incident'],
        ['co', 'http://nationalrail.co.uk/xml/common']
        ];
    iid  bigint; -- internal ID for incident
    inum char(32); -- incident number
    iver bigint;
    axml xml;
    bxml xml;
    aref char(2); -- affected toc ref
    rec  record;
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

    raise notice 'got %v', (xpath('//ic:ParticipantRef/text()', axml, ns))[1]::text;

    -- The main entry
    insert into darwin.incident (incidentnumber, participantref, created, summary, validfrom, validto, version)
    values (inum,
            (xpath('//ic:ParticipantRef/text()', axml, ns))[1]::text,
            (xpath('//ic:CreationTime/text()', axml, ns))[1]::text::timestamp with time zone,
            (xpath('//ic:Summary/text()', axml, ns))[1]::text,
            (xpath('//ic:ValidityPeriod/co:StartTime/text()', axml, ns))[1]::text::timestamp with time zone,
            (xpath('//ic:ValidityPeriod/co:EndTime/text()', axml, ns))[1]::text::timestamp with time zone,
            iver)
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
            insert into darwin.incident_affects_toc (ref, name)
            values ((xpath('//ic:OperatorRef/text()', bxml, ns))[1]::text,
                    (xpath('//ic:OperatorName/text()', bxml, ns))[1]::text)
            on conflict do nothing
            returning ref into aref;

            insert into darwin.incident_affects (id, toc)
            values (iid, aref)
            on conflict do nothing;
        end loop;
end ;
$$ language plpgsql;
