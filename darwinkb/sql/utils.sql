-- ======================================================================
-- text_to_jsonb converts a text into jsonb
--
-- Unlike to_jsonb() this will ensure that:
-- "true" or "false" are represented as their correct boolean types
-- Numbers are also in their correct json type
-- ======================================================================
create or replace function text_to_json(pval text)
    returns jsonb
as
$$
begin
    case
        when pval is null then
            return null;
        when pval = 'true' then
            return to_jsonb(true);
        when pval = 'false' then
            return to_jsonb(false);
        when pval ~ '^(-)?[0-9]+(\.[0-9]+)?$' then
            begin
                return to_jsonb(pval::numeric);
            exception
                when numeric_value_out_of_range then
                -- ignore as it's too big to fit in a real
            end; else
        -- ignore, keep value as is
        end case;
    return to_jsonb(pval);
end;
$$
    language plpgsql;

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

    -- If entry already exists then ensure it's an array and append the value to it
    if result ? pkey then
        ary = result -> pkey;

        if jsonb_typeof(ary) != 'array' then
            ary = jsonb_build_array(ary);
        end if;

        result = jsonb_set(result, array [pkey], ary || pval);
    else
        -- it doesn't exist so the value will be a single value
        result = jsonb_set(result, array [pkey], pval, true);
    end if;

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
    for rec in select '@' || (xpath('local-name(' || root || '/@*[' || i || '])', p_xml))[1]::text as rk,
                      v::text                                                                      as v
               from unnest(xpath(root || '/@*', p_xml)) with ordinality as a(v, i)
               where v is not null
               order by i
        loop
            result = jsonb_insert_upgrade(result, rec.rk, text_to_json(rec.v));
        end loop;

    -- children
    for rec in select (xpath('local-name(' || root || '/*[' || i || '])', p_xml))[1]::text as rk, v, i
               from unnest(xpath(root || '/*', p_xml)) with ordinality as a(v, i)
               order by i
        loop
            child = xml_to_json(rec.v);
            if child is not null then
                result = jsonb_insert_upgrade(result, rec.rk, child);
            end if;
        end loop;

    -- Add inner text as valueroot element .text(). No text (or just whitespace) will not include a value
    txt = regexp_replace(
            regexp_replace(array_to_string((xpath('/' || root || '/text()', p_xml))::text[], ' ', ''), '^\s+', ''),
            '\s+$', '');
    if txt is not null and txt != '' then
        -- Handle case of element containing just _value (no children or attrs)
        select into size count(*) from jsonb_object_keys(result);
        if size = 0 then
            result = text_to_json(txt);
        else
            result = jsonb_insert_upgrade(result, '_value', text_to_json(txt));
        end if;
    end if;

    return result;
end
$$ language plpgsql immutable;
