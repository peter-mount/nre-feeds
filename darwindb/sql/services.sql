-- services.sql handles basic queries against the service table
drop function darwin.getservices(pcrs char(3), pts timestamp with time zone);

create or replace function darwin.getservices(pcrs char(3), pts timestamp with time zone)
    returns table
            (
                tiploc varchar(7),
                rid bigint,
                location json,
                destination varchar(26),
                cancelled boolean,
                cancreason text,
                delayreason text,
                uid text,
                status text,
                trainid text,
                passengerservice bool,
                charterservice bool,
                toc text
            )
as
$$
declare
    fts timestamp without time zone = date_trunc('hour',pts);
    begin return query
select s.tiploc,
       s.rid,
       -- The JSON for this entry
       (select j.*
        from json_array_elements(sh.data -> 'locations') j
        where j ->> 'tiploc' = s.tiploc
          and (j -> 'timetable' ->> 'time')::time = s.ts::time
          and j ->> 'type' = s.type
       ),
       -- destination as text, falsedest overrides it if provided
       case
           when tf.name is not null then tf.name
           when td.name is not null then td.name
           when s.falsedest is not null then s.falsedest
           else s.destination
           end,
       s.cancelled,
       cr.cancel,
       dr.late,
       sh.data ->> 'uid',
       sh.data ->> 'status',
       sh.data ->> 'trainId',
       (sh.data ->> 'passengerService')::bool is not null,
       (sh.data ->> 'charterService')::bool is not null,
       sh.data ->> 'toc'
from darwin.service s
         inner join naptan.railreferences t on s.tiploc = t.tiploccode
         inner join darwin.schedule sh on s.rid = sh.rid
         left outer join timetable.tiploc td on s.destination = td.tiploc
         left outer join timetable.tiploc tf on s.falsedest = tf.tiploc
         left outer join darwin.reason cr on cr.id = s.cancreason
         left outer join darwin.reason dr on dr.id = s.delayreason
where t.crscode = pcrs
  and s.ts between fts and fts + '1 hour'::interval
order by ts;

end;
$$
    language plpgsql;