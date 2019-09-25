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
    aval   jsonb;
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

    aval = pval;
    if jsonb_typeof(aval) = 'string' then

        case
            when aval::text = '"true"' then
                aval = to_jsonb(true);
            when aval::text = '"false""' then
                aval = to_jsonb(false);
            when aval::text ~ '^"(-)?[0-9]+"$' then
                begin
                    aval = to_jsonb(replace(aval::text, '"', '')::bigint);
                exception
                    when numeric_value_out_of_range then
                    -- ignore as it's too big to fit in a bigint
                end; else
            -- ignore, keep value as is
            end case;
    end if;

    -- If entry already exists then ensure it's an array and append the value to it
    if result ? pkey then
        ary = result -> pkey;

        if jsonb_typeof(ary) != 'array' then
            ary = jsonb_build_array(ary);
        end if;

        ary = jsonb_insert(ary, array ['0'], aval, true);

        -- remove the existing entry as we will insert it with the new value shortly
        result = result - pkey;
    else
        -- it doesn't exist so the value will be a single value
        ary = aval;
    end if;

    -- Set the new value, either the raw value or the updated array
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
    txt    text;
    rec    record;
    child  jsonb;
begin

    -- Get root element name
    select '/*[name() = "' || (xpath('name(/*)', p_xml))[1]::text || '"]' into root;

    -- attributes
    for rec in
        select '@' || (xpath('local-name(' || root || '/@*[' || i || '])', p_xml))[1]::text as rk,
               to_jsonb(v::text)                                                            as v
        from unnest(xpath(root || '/@*', p_xml)) with ordinality as a(v, i)
        where v is not null
        order by i
        loop
            result = jsonb_insert_upgrade(result, rec.rk, rec.v);
        end loop;

    -- children
    for rec in select (xpath('local-name(' || root || '/*[' || i || '])', p_xml))[1]::text as rk, v
               from unnest(xpath(root || '/*', p_xml)) with ordinality as a(v, i)
               order by i
        loop
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
        if size = 0 then
            result = to_jsonb(txt);
        else
            result = jsonb_insert_upgrade(result, '_value', to_jsonb(txt));
        end if;
    end if;

    return result;
end
$$ language plpgsql immutable;
