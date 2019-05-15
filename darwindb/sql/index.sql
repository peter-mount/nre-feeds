-- index.sql Manages the indexing of schedules by station
--
-- This is used by uktra.in to be able to search for historical services
--

-- service entries by tiploc, rid and ts (time at a tiploc)
-- Note this only contains those scheduled to arrive or depart at the tiploc
--drop table darwin.service cascade;

create table if not exists darwin.service
(
    tiploc      varchar(7)                  not null,
    rid         bigint                      not null,
    ts          timestamp without time zone not null,
    type        varchar(4)                  not null,
    pta         time,
    ptd         time,
    plat        varchar(16),
    arrive      time,
    depart      time,
    activity    text,
    destination varchar(7)                  not null,
    falsedest   varchar(7),
    cancelled   boolean                     not null default false,
    cancreason  int                         not null default 0,
    delayreason int                         not null default 0
) partition by range (ts);

create unique index if not exists service_ttr on darwin.service (ts, tiploc, rid);
create index if not exists service_tt on darwin.service (ts, tiploc);
create index if not exists service_r on darwin.service (rid);


-- Called by addservice when inserting an entry into the darwin.service table to create a partition
-- if one doesn't exist.
-- There is one partition per date (initially)
create or replace function darwin.serviceCreateTable(pts timestamp with time zone)
    returns void
as
$$
declare
    sts timestamp with time zone = date_trunc('day', pts);
    ets timestamp with time zone = sts + '1 day'::interval ;
    s   text;
begin
    execute format(
            'create table if not exists darwin.service_%s partition of darwin.service for values from (%L) to (%L)',
            to_char(sts, 'YYYYMMDD'),
            sts,
            ets
        );
end
$$
    language plpgsql;

-- Add a single entry into the service table
create or replace function darwin.addservice(prow darwin.service)
    returns void
as
$$
declare
begin
    insert into darwin.service (tiploc, rid, ts, pta, ptd, plat, arrive, depart, destination, falsedest, type,
                                cancelled, cancreason,
                                delayreason, activity)
    values (prow.tiploc, prow.rid, prow.ts, prow.pta, prow.ptd, prow.plat, prow.arrive, prow.depart, prow.destination,
            prow.falsedest, prow.type,
            prow.cancelled, prow.cancreason, prow.delayreason, prow.activity);
exception
    when unique_violation then
        update darwin.service
        set falsedest   = prow.falsedest,
            type        = prow.type,
            cancelled   = prow.cancelled,
            cancreason  = prow.cancreason,
            delayreason = prow.delayreason,
            activity    = prow.activity,
            plat        = prow.plat,
            arrive      = prow.arrive,
            depart      = prow.depart
        where tiploc = prow.tiploc
          and rid = prow.rid
          and ts = prow.ts;

    when check_violation then
        execute darwin.serviceCreateTable(prow.ts);

        begin
            insert into darwin.service (tiploc, rid, ts, pta, ptd, plat, arrive, depart, destination, falsedest, type,
                                        cancelled, cancreason,
                                        delayreason, activity)
            values (prow.tiploc, prow.rid, prow.ts, prow.pta, prow.ptd, prow.plat, prow.arrive, prow.depart,
                    prow.destination,
                    prow.falsedest, prow.type,
                    prow.cancelled, prow.cancreason, prow.delayreason, prow.activity);
        exception
            when unique_violation then
                update darwin.service
                set falsedest   = prow.falsedest,
                    type        = prow.type,
                    cancelled   = prow.cancelled,
                    cancreason  = prow.cancreason,
                    delayreason = prow.delayreason,
                    activity    = prow.activity,
                    plat        = prow.plat,
                    arrive      = prow.arrive,
                    depart      = prow.depart
                where tiploc = prow.tiploc
                  and rid = prow.rid
                  and ts = prow.ts;
        end;
end;
$$
    language plpgsql;

-- Index a single service
create or replace function darwin.indexservice(prid bigint)
    returns void
as
$$
declare
    i     json;
    j     json;
    ts    timestamp with time zone;
    first boolean = false;
    fts   timestamp with time zone;
    prow  darwin.service;
begin

    select into j data from darwin.schedule where rid = prid;
    if not found then
        return;
    end if;

    ts = (j ->> 'ssd')::date::timestamp with time zone;

    prow.rid = (j ->> 'rid')::bigint;
    prow.cancreason = (j -> 'cancelReason' ->> 'reason')::int;
    prow.delayreason = (j -> 'lateReason' ->> 'reason')::int;
    prow.destination = j -> 'destinationLocation' ->> 'tiploc';

    -- Some deactivations come with no destination
    if prow.destination is null then
        prow.destination = 'NULLTPL';
    end if;

    for i in select * from json_array_elements(j -> 'locations')
        loop
            if i -> 'timetable' is not null
                and i -> 'forecast' is not null
                and i ->> 'type' IN ('OR', 'OPOR', 'IP', 'DT', 'OPDT')
            then

                prow.type = i ->> 'type';
                prow.pta = (i -> 'timetable' ->> 'pta')::time;
                prow.ptd = (i -> 'timetable' ->> 'ptd')::time;
                -- The display time as a timestamp & if not the first entry and before the first entry time then presume
                -- we have crossed midnight.
                prow.ts = ts + (i -> 'timetable' ->> 'time')::TIME;
                if first then
                    if prow.ts < fts then
                        prow.ts = prow.ts + '1 day'::interval;
                    end if;
                else
                    first = true;
                    fts = prow.ts;
                end if;

                prow.plat = i -> 'forecast' -> 'plat' ->> 'plat';

                -- The recorded at times for arr & dep. Don't use et else it appears the service has run in the future!
                prow.arrive = (i -> 'forecast' -> 'arr' ->> 'at')::TIME;
                prow.depart = (i -> 'forecast' -> 'dep' ->> 'at')::TIME;

                prow.activity = i -> 'planned' ->> 'activity';

                prow.tiploc = i ->> 'tiploc';
                prow.falsedest = i ->> 'falsedest';

                prow.cancelled = i ->> 'canc' is not null;

                execute darwin.addservice(prow);
            end if;
        end loop;
end;
$$
    language plpgsql;

-- table to capture errors from indexservices()
create table if not exists darwin.indexerrors
(
    rid      bigint not null,
    ts       timestamp with time zone,
    failtime timestamp with time zone,
    msg      text,
    detail   text,
    hint     text,
    context  text
);
truncate darwin.indexerrors;

-- Indexes the first 500 entries in the scheduleupdate table.
-- It takes them in order, oldest first.
create or replace function darwin.indexservices()
    returns int
as
$$
declare
    rec        record;
    msgText    text;
    pgexdetail text;
    pgexhint   text;
    pgexctx    text;
    ctr        int = 0;
begin

    for rec in select rid, date
               from darwin.scheduleupdate
               where date is null
                  or date < now() - '10 minutes'::interval
               order by date
               limit 500
        loop
            begin
                execute darwin.indexservice(rec.rid);

                if rec.date is null then
                    delete
                    from darwin.scheduleupdate
                    where rid = rec.rid
                      and date is null;
                else
                    delete
                    from darwin.scheduleupdate
                    where rid = rec.rid
                      and date = rec.date;
                end if;

                ctr = ctr + 1;
            exception
                when others then
                    get stacked diagnostics msgText = MESSAGE_TEXT ,
                        pgexdetail = PG_EXCEPTION_DETAIL,
                        pgexhint = PG_EXCEPTION_HINT,
                        pgexctx = PG_EXCEPTION_CONTEXT;
                    insert into darwin.indexerrors
                    values (rec.rid, rec.date, now(), msgText, pgexdetail, pgexhint, pgexctx);

                    ctr = 0;
            end;
        end loop;
    return ctr;
end;
$$
    language plpgsql;
