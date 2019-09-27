-- ======================================================================
-- station.sql handles the station reference feed consisting of details
-- about the station, facilities etc as well as being a postgis layer
-- ======================================================================
create schema if not exists darwin;

create table if not exists darwin.station
(
    crs             char(3)                  not null,              -- CRS code
    name            name                     not null,              -- station name
    name16          varchar(16)              not null,              -- 16 character name
    longitude       real                     not null,              -- longitude
    latitude        real                     not null,              -- latitude,
    stationOperator char(2) references darwin.companies (atoccode), -- Station operator
    updated         timestamp with time zone not null,              -- Timestamp entry was last updated
    updatedBy       name,                                           -- Who last updated the entry
    primary key (crs)
);
select addgeometrycolumn('', 'darwin', 'station', 'geom', 4326, 'POINT', 2, true);

create table if not exists darwin.station_xml
(
    crs     char(3)                  not null references darwin.station (crs), -- fk to station table
    updated timestamp with time zone not null,                                 -- Timestamp entry was last updated
    xml     xml                      not null,                                 -- original xml
    json    jsonb                    not null,                                 -- generated json
    primary key (crs)
);

-- ======================================================================
create or replace function darwin.updateStation(pxml xml)
    returns void
as
$$
declare
    ns       text[][] := ARRAY [
        ['add','http://www.govtalk.gov.uk/people/AddressAndPersonalDetails'],
        ['com','http://nationalrail.co.uk/xml/common'],
        ['stn','http://nationalrail.co.uk/xml/station']
        ];
    acrs     char(3);
    axml     xml;
    along    real;
    alat     real;
    aupdated timestamp with time zone;
begin
    foreach axml in array (xpath('//stn:StationList/stn:Station', pxml, ns))
        loop
            acrs = (xpath('//stn:Station/stn:CrsCode/text()', axml, ns))[1]::text;
            along = (xpath('//stn:Station/stn:Longitude/text()', axml, ns))[1]::text::real;
            alat = (xpath('//stn:Station/stn:Latitude/text()', axml, ns))[1]::text::real;
            aupdated = (xpath('//stn:Station/stn:ChangeHistory/com:LastChangedDate/text()', axml,
                              ns))[1]::text::timestamp with time zone;
            -- Add main entry, updating only if changed
            insert into darwin.station as s
            (crs,
             name, name16,
             longitude, latitude,
             stationoperator,
             updated, updatedBy,
             geom)
            values (acrs,
                    (xpath('//stn:Station/stn:Name/text()', axml, ns))[1]::text,
                    (xpath('//stn:Station/stn:SixteenCharacterName/text()', axml, ns))[1]::text,
                    along,
                    alat,
                    (xpath('//stn:Station/stn:StationOperator/text()', axml, ns))[1]::text,
                    aupdated,
                    (xpath('//stn:Station/stn:ChangeHistory/com:ChangedBy/text()', axml, ns))[1]::text,
                    case
                        when along is not null and alat is not null then
                            public.ST_SetSRID(public.ST_MakePoint(along::double precision, alat::double precision),
                                              4326)
                        else null
                        end)
            on conflict (crs)
                do UPDATE set name=excluded.name,
                              name16 = excluded.name16,
                              longitude = excluded.longitude,
                              latitude = excluded.latitude,
                              stationOperator = excluded.stationOperator,
                              updated = excluded.updated,
                              updatedBy = excluded.updatedBy
            where s.updated < excluded.updated;

            insert into darwin.station_xml as x
                (crs, updated, xml, json)
            values (acrs, aupdated, axml, public.xml_to_json(axml))
            on conflict (crs)
                do update set xml=excluded.xml,
                              json=excluded.json,
                              updated = excluded.updated
            where x.updated < excluded.updated;
        end loop;
end;
$$ language plpgsql;