-- ======================================================================
-- companies.sql handles the companies static knowledge base feed
-- ======================================================================
create schema if not exists darwin;

create table if not exists darwin.companies
(
    atocCode        char(2)                  not null, -- ATOC code
    name            name                     not null, -- public name
    legalName       name                     not null, -- legal name
    atocMember      bool                     not null, -- ATOC member
    stationOperator bool                     not null, -- Operates stations
    operatingFrom   timestamp with time zone not null, -- start date of operations
    operatingTo     timestamp with time zone,          -- end date, null if unknown
    deleted         bool                     not null, -- set to true if entry nolonger in feed - used to keep referential integrity
    primary key (atocCode)
);

create table if not exists darwin.companies_xml
(
    atocCode char(2) not null references darwin.companies (atocCode), -- fk to companies table
    xml      xml     not null,                                        -- original xml
    json     jsonb   not null,                                        -- generated json
    primary key (atocCode)
);

-- ======================================================================
create or replace function darwin.updateCompanies(pxml xml)
    returns void
as
$$
declare
    ns   text[][] := ARRAY [
        ['add','http://www.govtalk.gov.uk/people/AddressAndPersonalDetails'],
        ['com','http://nationalrail.co.uk/xml/common'],
        ['toc','http://nationalrail.co.uk/xml/toc']
        ];
    atc  char(2);
    axml xml;
begin
    -- Mark all entries as deleted
    update darwin.companies set deleted= true where not deleted;

    -- Recreate table content
    foreach axml in array (xpath('//toc:TrainOperatingCompanyList/toc:TrainOperatingCompany', pxml, ns))
        loop
            atc = (xpath('//toc:TrainOperatingCompany/toc:AtocCode/text()', axml, ns))[1]::text;

            -- Add main entry
            insert into darwin.companies (atocCode, name, legalName, atocMember, stationOperator, operatingFrom,
                                          operatingTo, deleted)
            values (atc,
                    (xpath('//toc:TrainOperatingCompany/toc:Name/text()', axml, ns))[1]::text,
                    (xpath('//toc:TrainOperatingCompany/toc:LegalName/text()', axml, ns))[1]::text,
                    (xpath('//toc:TrainOperatingCompany/toc:AtocMember/text()', axml, ns))[1]::text = 'true',
                    (xpath('//toc:TrainOperatingCompany/toc:StationOperator/text()', axml, ns))[1]::text = 'true',
                    (xpath('//toc:TrainOperatingCompany/toc:OperatingPeriod/toc:StartDate/text()', axml,
                           ns))[1]::text::timestamp with time zone,
                    (xpath('//toc:TrainOperatingCompany/toc:OperatingPeriod/toc:EndDate/text()', axml,
                           ns))[1]::text::timestamp with time zone,
                    false)
            on conflict (atocCode)
                do UPDATE set name=excluded.name,
                              legalName=excluded.legalName,
                              atocMember=excluded.atocMember,
                              stationOperator=excluded.stationOperator,
                              operatingFrom=excluded.operatingFrom,
                              operatingTo=excluded.operatingTo,
                              deleted= false;

            insert into darwin.companies_xml (atoccode, xml, json)
            values (atc, axml, xml_to_json(axml))
            on conflict (atocCode)
                do update set xml=excluded.xml, json=excluded.json;
        end loop;
end;
$$ language plpgsql;