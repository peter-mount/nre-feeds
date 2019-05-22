-- performance.sql contains a procedure to return the performance of a service over time.
--

-- Helper, calculates the delay as an interval.
--
-- To get the delay in seconds use: extract(epoch from darwin.delay(act,exp))
-- To get the delay in minutes use: extract(epoch from darwin.delay(act,exp))/60.0
create or replace function darwin.delay(actual time, expected time)
    returns interval
as
$$
select case
           when (actual - expected) > '18 hours' then expected - actual + '24:00'
           when (actual - expected) < '-18 hours' then actual - expected + '24:00'
           else actual - expected
           end
$$
    language sql;

-- Initially we will presume that the service's UID is the unique key
create or replace function darwin.uidperformance(puid text, pcrs text)
    returns json
as
$$
with tiplocs as (select tiploc from timetable.tiploc where crs = pcrs),
     entries as (
         select s.*
         from darwin.service s
         where s.uid = puid
           and s.tiploc in ((select tiploc from tiplocs))
         order by s.ts
     ),
     services as (
         select rid, data
         from darwin.schedule
         where rid in (select rid from entries)
     )
select json_build_object(
           -- report name
               'name', 'uidperformance',
           -- date generated
               'generated', now(),
           -- last modified timestamp = latest date in data
               'lastModified', (select max((data ->> 'date')::timestamp with time zone) from services),
           -- report parameters
               'params', json_build_object(
                       'uid', puid,
                       'crs', pcrs
                   ),
           -- report ranges
               'range', json_build_object(
                       'delay',
                       json_build_object(
                               'min', (select min(least(extract(epoch from darwin.delay(arrive, pta)) / 60.0,
                                                        extract(epoch from darwin.delay(depart, ptd)) / 60.0))
                                       from entries),
                               'max', (select max(greatest(extract(epoch from darwin.delay(arrive, pta)) / 60.0,
                                                           extract(epoch from darwin.delay(depart, ptd)) / 60.0))
                                       from entries)
                           ),
                       'ts',
                       json_build_object(
                               'min', (select min(ts) from entries),
                               'max', (select max(ts) from entries)
                           )
                   ),
           -- report data
               'data', (select json_agg(entry)
                        from (select json_build_object(
                                             'rid', e.rid,
                                             'ts', e.ts,
                                             'location', l,
                                             'delay', json_build_object(
                                                     'arr', extract(epoch from darwin.delay(e.arrive, e.pta)) / 60.0,
                                                     'dep', extract(epoch from darwin.delay(e.depart, e.ptd)) / 60.0
                                                 ),
                                             'reason', case
                                                           when (l ->> 'cancelled')::boolean and
                                                                (s.data -> 'cancelReason' ->> 'reason')::int > 0
                                                               then s.data -> 'cancelReason'
                                                           when (s.data -> 'lateReason' ->> 'reason')::int > 0
                                                               then s.data -> 'lateReason'
                                                           else null
                                                 end,
                                             'destination', s.data -> 'destinationLocation' ->> 'tiploc'
                                         ) as entry
                              from entries e
                                       inner join services s on e.rid = s.rid
                                       inner join json_array_elements(s.data -> 'locations') l
                                                  on l ->> 'tiploc' in ((select tiploc from tiplocs))
                                                      and (l -> 'timetable' ->> 'time')::time = e.ts::time
                                                      and l ->> 'type' = e.type
                              order by e.ts) t),
           -- cancel/late reason xref
               'reason', (select json_object_agg(id, obj)
                          from (
                                   select distinct on (id) r.id, row_to_json(r.*) as obj
                                   from darwin.reason r
                                   where r.id in (
                                       select distinct (cs.data -> 'cancelReason' ->> 'reason')::int
                                       from services cs
                                       union
                                       select distinct (ls.data -> 'lateReason' ->> 'reason')::int
                                       from services ls
                                   )
                               ) t),
           -- tiploc xref
               'tiploc', (select json_object_agg(tiploc, obj)
                          from (
                                   with s as (select json_array_elements(data -> 'locations') as loc from services)
                                   select distinct on ( tiploc ) tpl.tiploc, row_to_json(tpl.*) as obj
                                   from timetable.tiploc tpl
                                   where tpl.tiploc in (
                                       -- All tiplocs for this crs
                                       select tiploc
                                       from tiplocs
                                            -- All destination tiplocs
                                       union
                                       select distinct s.data -> 'destinationLocation' ->> 'tiploc'
                                       from services s
                                            -- Any falseDestination tiplocs
                                       union
                                       select distinct loc ->> 'falseDestination'
                                       from s
                                       where loc ->> 'tiploc' in ((select tiploc from tiplocs))
                                             -- cancelReason at/near tiplocs
                                       union
                                       select distinct sv.data -> 'cancelReason' ->> 'tiploc'
                                       from services sv
                                            -- lateReason at/near tiplocs
                                       union
                                       select distinct sv.data -> 'lateReason' ->> 'tiploc'
                                       from services sv
                                   )
                               ) t)
           );
$$
    language sql;
