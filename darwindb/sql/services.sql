-- services.sql handles basic queries against the service table
drop function darwin.getservices(pcrs char(3), pts timestamp with time zone);

create or replace function darwin.getservices(pcrs char(3), pts timestamp with time zone)
    returns table
            (
                tiploc varchar(7),
                rid bigint,
                ts timestamp without time zone,
                type varchar(4),
                pta time,
                ptd time,
                plat varchar(16),
                arrive time,
                depart time,
                activity text,
                destination varchar(26),
                falsedest varchar(26),
                cancelled boolean,
                cancreason int,
                delayreason int,
                uid text,
                status text,
                trainid text,
                passengerservice bool,
                charterservice bool
            )
as
$$
declare
    fts timestamp without time zone = date_trunc('hour',pts);
    begin return query
select s.tiploc,
       s.rid,
       s.ts,
       s.type,
       s.pta,
       s.ptd,
       s.plat,
       s.arrive,
       s.depart,
       s.activity,
       --s.destination,
       case
           when td.name is not null then td.name
           else s.destination
           end,
       --s.falsedest,
       case
           when tf.name is not null then tf.name
           else s.falsedest
           end,
       s.cancelled,
       s.cancreason,
       s.delayreason,
       sh.data ->> 'uid'                                  as uid,
       sh.data ->> 'status'                               as status,
       sh.data ->> 'trainId'                              as trainid,
       (sh.data ->> 'passengerService')::bool is not null as passengerService,
       (sh.data ->> 'charterService')::bool is not null   as charterService
from darwin.service s
         inner join naptan.railreferences t on s.tiploc = t.tiploccode
         inner join darwin.schedule sh on s.rid = sh.rid
         left outer join timetable.tiploc td on s.destination = td.tiploc
         left outer join timetable.tiploc tf on s.falsedest = tf.tiploc
where t.crscode = pcrs
  and s.ts between fts and fts + '1 hour'::interval
order by ts;

end;
$$
    language plpgsql;