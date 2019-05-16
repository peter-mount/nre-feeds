-- reference.sql An implementation of the reference schema in postgresql
--
-- Note, initially darwinkb will continue to use bolt as thats ideal, but we have this
-- so that we can keep referential integrity & future reuse

-- A copy of the Cancel/Delay reasons
create table darwin.reason
(
    id     int not null,
    cancel text,
    late   text,
    primary key (id)
);

create or replace function darwin.updateReasons(pjson json)
    returns void
as
$$
declare
    reason json;
begin
    for reason in select json_array_elements(pjson -> 'reasons')
        loop
            if (reason ->> 'canc')::boolean then
                insert into darwin.reason (id, cancel)
                values ((reason ->> 'code')::int, reason ->> 'reasontext')
                on conflict (id)
                    do update
                    set cancel = EXCLUDED.cancel;
            else
                insert into darwin.reason (id, late)
                values ((reason ->> 'code')::int, reason ->> 'reasontext')
                on conflict(id)
                    do update
                    set late = EXCLUDED.late;
            end if;
        end loop;

end;
$$
    language plpgsql;