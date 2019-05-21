-- services.sql handles basic queries against the service table
drop function darwin.getservices(pcrs char(3), pts timestamp with time zone);

create or replace function darwin.getservices(pcrs char(3), pts timestamp with time zone)
    returns json
as
$$
declare
    -- start & end time to scan. This is inclusive hence not adding 1 hour to get ets
    fts       timestamp without time zone = date_trunc('hour', pts);
    ets       timestamp without time zone = fts + '59 minutes 59 seconds'::interval;
    -- json results
    jstation  json;
    jservices json;
    jreason   json;
    jtiplocs  json;
begin

    -- Station details, first entry should be official name
    select into jstation json_agg(tpl)
    from (select st.tiploc as tpl
          from timetable.tiploc st
          where st.crs = pcrs
          order by st.nlcdesc desc
         ) t;

    -- Get the available services
    select into jservices json_agg(service)
    from (select json_build_object(
                         'rid', s.rid,
                     -- The JSON for this entry
                         'location', (select j.*
                                      from json_array_elements(sh.data -> 'locations') j
                                      where j ->> 'tiploc' = s.tiploc
                                        and (j -> 'timetable' ->> 'time')::time = s.ts::time
                                        and j ->> 'type' = s.type
                                      limit 1
                         ),
                     -- destination, falsedest overrides it if provided
                         'destination', case
                                            when s.falsedest is not null then s.falsedest
                                            when s.destination is not null then s.destination
                                            else ''
                             end,
                         'cancelled', s.cancelled,
                         'cancelReason', sh.data -> 'cancelReason',
                         'lateReason', sh.data -> 'lateReason',
                         'uid', sh.data ->> 'uid',
                         'status', sh.data ->> 'status',
                         'trainId', sh.data ->> 'trainId',
                         'passengerService', (sh.data ->> 'passengerService')::bool is not null,
                         'charterService', (sh.data ->> 'charterService')::bool is not null,
                         'toc', sh.data ->> 'toc',
                         'association', sh.data -> 'association',
                         'formation', sh.data -> 'formation'
                     ) as service
          from darwin.service s
                   inner join timetable.tiploc t on s.tiploc = t.tiploc
                   inner join darwin.schedule sh on s.rid = sh.rid
                   left outer join darwin.reason cr on cr.id = s.cancreason
                   left outer join darwin.reason dr on dr.id = s.delayreason
          where t.crs = pcrs
            and s.ts between fts and ets
          order by ts
         ) t;

    -- resolve reasons
    select into jreason json_object_agg(id, obj)
    from (select distinct on (r.id) r.id, row_to_json(r.*) as obj
          from darwin.reason r
          where r.id in (
              select distinct (value -> 'cancelReason' ->> 'reason')::int
              from json_array_elements(jservices) cr
              union
              select distinct (value -> 'delayReason' ->> 'reason')::int
              from json_array_elements(jservices) dr
          )
         ) t;

    -- resolve tiplocs
    select into jtiplocs json_object_agg(tiploc, obj)
    from (
             select distinct on (tiploc) tpl.tiploc, row_to_json(tpl.*) as obj
             from timetable.tiploc tpl
             where tpl.tiploc in (
                 -- station tiploc entries
                 select distinct *
                 from json_array_elements_text(jstation) l
                      -- destination entries
                 union
                 select distinct s ->> 'destination'
                 from json_array_elements(jservices) s
                      -- reason entries
                 union
                 select distinct s -> 'cancelReason' ->> 'tiploc'
                 from json_array_elements(jservices) s
                 union
                 select distinct s -> 'delayReason' ->> 'tiploc'
                 from json_array_elements(jservices) s
                      -- schedule locations
                 union
                 select distinct s -> 'location' ->> 'tiploc'
                 from json_array_elements(jservices) s
                      -- association locations, where a service has them
                 union
                 select distinct b ->> 'tiploc'
                 from json_array_elements(jservices) s
                          left outer join json_array_elements(s -> 'association') b on true
                 where s ->> 'association' != 'null'
             )
         ) t;

    -- final result
    return json_build_object(
            'station', jstation,
            'services', jservices,
            'reason', jreason,
            'tiploc', jtiplocs
        );

end ;
$$
    language plpgsql;

-- Retrieves all details about a service
create or replace function darwin.getservice(prid bigint)
    returns json
as
$$
declare
    schedule json;
    tiplocs  json;
begin
    select into schedule data from darwin.schedule where rid = prid;
    if not found then
        return null;
    end if;

    select into tiplocs json_object_agg(tiploc, obj)
    from (
             select tiploc,
                    json_build_object(
                            'tiploc', tiploc,
                            'crs', trim(crs),
                            'name', name
                        ) as obj
             from timetable.tiploc
             where tiploc in (
                 select distinct l ->> 'tiploc'
                 from json_array_elements(schedule -> 'locations') l
             )) t;

    return json_build_object(
            'rid', prid,
            'schedule', schedule,
            'tiploc', tiplocs
        );
end ;
$$
    language plpgsql;