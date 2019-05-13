-- schedule.sql - the main active schedule

-- Just a store of the json
create table if not exists darwin.schedule
(
    rid  bigint not null,
    data text   not null,
    date timestamp with time zone,
    primary key (rid)
) partition by range (rid);
create index if not exists schedule_dt on darwin.schedule (date);

-- Called by updateschedule when it attempts to insert into an unknown partition.
-- This will create a partition schedule_YYYYMMDD where YYYYMM is the month and DD is the date at the
-- begining of the partition area defined by daysInPartition, default is 4 so we have a partition for blocks
-- of 4 days due to the number of services involved.
create or replace function darwin.scheduleCreateTable(prid bigint)
    returns void
as
$$
declare
    daysInPartition bigint = 4;
    baseidx         bigint = prid / 10000000;
    monthidx        bigint = baseidx / 100;
    weekidx         bigint = ((baseidx % 100) / daysInPartition);
    fromrid         bigint = ((monthidx * 100) + (daysInPartition * weekidx)) * 10000000;
    torid           bigint = fromrid + (daysInPartition * 10000000);
    tname           text   = concat('schedule_', monthidx, '_', weekidx);
    s               text;
begin
    s = format(
            'create table if not exists darwin.%I partition of darwin.schedule for values from (%s) to (%s)',
            tname,
            fromrid,
            torid
        );
    execute s;
end
$$
    language plpgsql;

-- Insert/Update a schedule entry
-- drop function darwin.updateschedule(json)
create or replace function darwin.updateschedule(pmsg json)
    returns void
as
$$
declare
    prid      bigint = (pmsg ->> 'RID')::bigint;
    pschedule json   = pmsg -> 'Schedule';
begin
    raise notice 'rid %', prid;
    insert into darwin.schedule (rid, data, date)
    values (prid,
            pschedule,
            now());
exception
    when unique_violation then
        raise notice 'update rid %', prid;
        update darwin.schedule
        set data = data,
            date = now()
        where rid = prid;

    when check_violation then
        execute darwin.scheduleCreateTable(prid);

        insert into darwin.schedule (rid, data, date)
        values (prid,
                pschedule,
                now());
end;
$$
    language plpgsql;

create or replace function darwin.getschedule(prid varchar(16))
    returns text
as
$$
select data
from darwin.schedule
where rid = prid;
$$
    language sql;