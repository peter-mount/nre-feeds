-- schedule.sql - the main active schedule

-- Just a store of the json
create table if not exists darwin.schedule
(
    rid  bigint not null,
    data json   not null,
    date timestamp with time zone,
    primary key (rid)
) partition by range (rid);
create index if not exists schedule_dt on darwin.schedule (date);
create index if not exists schedule_rdt on darwin.schedule (rid, date);

-- Table holding RID's that need reindexing
create table if not exists darwin.scheduleUpdate
(
    rid  bigint not null,
    date timestamp with time zone,
    primary key (rid)
);

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
--
-- Note: we can't use an upset (insert on conflict(rid) as it doesn't
--       work fully with partitions just yet, so if we get a:
--          unique violation then do an update
--          check_violation then ensure we have the partition table,
--              then do the insert/update again.
--
-- Also we check the schedule date so if we do an update then only do it
-- if the schedule date is later than the one in the db - to reduce
-- the number of updates, especially if we are doing a resync!
--
create or replace function darwin.updateschedule(pmsg json)
    returns void
as
$$
declare
    prid      bigint                   = (pmsg ->> 'RID')::bigint;
    pschedule json                     = pmsg -> 'Schedule';
    pdate     timestamp with time zone = (pschedule ->> 'date')::timestamp with time zone;
begin

    -- Ensure we have a date, don't know why we can have a null date inbound
    if pdate is null then
        pdate = now();
    end if;

    insert into darwin.scheduleUpdate (rid, date)
    values (prid, pdate)
    on conflict (rid)
        do update
        set date = excluded.date;

    begin
        insert into darwin.schedule (rid, data, date)
        values (prid, pschedule, pdate);
    exception
        when unique_violation then
            update darwin.schedule
            set data = pschedule,
                date = pdate
            where rid = prid
              and date < pdate;

        when check_violation then
            execute darwin.scheduleCreateTable(prid);

            begin
                insert into darwin.schedule (rid, data, date)
                values (prid, pschedule, pdate);
            exception
                when unique_violation then
                    update darwin.schedule
                    set data = pschedule,
                        date = pdate
                    where rid = prid
                      and date < pdate;
            end;
    end;
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