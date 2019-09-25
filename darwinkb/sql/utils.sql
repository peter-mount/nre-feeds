-- ======================================================================
-- Adds an entry to a jsonb object similar to jsonb_insert
-- except existing keys are converted into an array
-- ======================================================================
create or replace function jsonb_insert_upgrade(ptarget jsonb, pkey text, pval jsonb)
    returns jsonb
as
$$
declare
    result jsonb;
    ary    jsonb;
begin
    -- Ensure we have an object
    result = ptarget;
    if result is null then
        result = '{}';
    end if;

    -- Ignore nulls
    if pval is null then
        return result;
    end if;

    if result ? pkey then
        ary = result -> pkey;

        if jsonb_typeof(ary) != 'array' then
            ary = jsonb_build_array(ary);
        end if;

        ary = jsonb_insert(ary, array ['0'], pval, true);

        result = result - pkey;
    else
        ary = pval;
    end if;

    raise notice '% = %', pkey, ary;
    result = jsonb_insert(result, array [pkey], ary, true);

    return result;
end;
$$
    language plpgsql;

-- ======================================================================
-- Converts XML into JSON
-- ======================================================================
--drop function xml_to_json(p_xml xml);
create or replace function xml_to_json(p_xml xml)
    returns jsonb
as
$$
declare
    result jsonb = '{}';
    root   text;
    size   int;
    keys   text[];
    vals   jsonb[];
    txt    text;
    rec    record;
    child  jsonb;
begin
    --raise notice 'xml %', p_xml;

    -- Get root element name
    select '/*[name() = "' || (xpath('name(/*)', p_xml))[1]::text || '"]' into root;
    raise notice 'root %', root;

    -- attributes
    for rec in
        select '@' || (xpath('local-name(' || root || '/@*[' || i || '])', p_xml))[1]::text as rk,
               to_jsonb(v::text)                                                            as v
        from unnest(xpath(root || '/@*', p_xml)) with ordinality as a(v, i)
        where v is not null
        order by i
        loop
            raise notice 'attr % %', rec.rk, rec.v;
            result = jsonb_insert_upgrade(result, rec.rk, rec.v);
        end loop;

    -- children
    for rec in select (xpath('local-name(' || root || '/*[' || i || '])', p_xml))[1]::text as rk, v
               from unnest(xpath(root || '/*', p_xml)) with ordinality as a(v, i)
               order by i
        loop
            raise notice 'child %', rec.rk;
            child = xml_to_json(rec.v);
            if child is not null then
                result = jsonb_insert_upgrade(result, rec.rk, child);
            end if;
        end loop;

    -- Add inner text as valueroot element .text(). No text (or just whitespace) will not include a value
    txt = regexp_replace(array_to_string((xpath('/' || root || '/text()', p_xml))::text[], ' ', ''), '\s+$', '');
    if txt is not null and txt != '' then
        -- Handle case of element containing just _value (no children or attrs)
        select into size count(*) from jsonb_object_keys(result);
        raise notice 'size %', size;
        if size = 0 then
            result = to_jsonb(txt);
        else
            result = jsonb_insert_upgrade(result, '_value', to_jsonb(txt));
        end if;
    end if;

    return result;
end
$$ language plpgsql immutable;
