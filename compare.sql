-- ferretdb_scalar_type returns v's type name and matching PostgreSQL value.
CREATE OR REPLACE FUNCTION ferretdb_scalar_type(v jsonb) RETURNS record AS $$
DECLARE
    t text;
    res record;
BEGIN
    t := jsonb_typeof(v);
    CASE t
        WHEN 'object' THEN
            CASE
                WHEN v->'$k' IS NOT NULL THEN
                    RAISE 'ferretdb_scalar_type: v is a document: %', v;
                WHEN v->'$f' IS NOT NULL THEN

                WHEN v->'$b' IS NOT NULL THEN
                WHEN v->'$o' IS NOT NULL THEN
                WHEN v->'$d' IS NOT NULL THEN
                WHEN v->'$r' IS NOT NULL THEN
                WHEN v->'$t' IS NOT NULL THEN
                WHEN v->'$l' IS NOT NULL THEN
                WHEN v->'$n' IS NOT NULL THEN
                WHEN v->'$c' IS NOT NULL THEN
                ELSE
                    RAISE 'ferretdb_scalar_type: v is an unexpected jsonb object: %', v;
            END CASE;
        WHEN 'array' THEN
            RAISE 'ferretdb_scalar_type: v is jsonb array: %', v;
        WHEN 'string' THEN
            -- use #>> operator to remove quotes, escaping, etc
            res := ('string'::text, v #>> '{}');
        WHEN 'number' THEN
            res := ('int32'::text, v::int4);
        WHEN 'boolean' THEN
            res := ('bool'::text, v::boolean);
        ELSE
            RAISE 'ferretdb_scalar_type: v has unhandled jsonb type: %', t;
    END CASE;

    RETURN res;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE OR REPLACE FUNCTION ferretdb_compare(a jsonb, b jsonb) RETURNS smallint AS $$
BEGIN

END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- CREATE OR REPLACE FUNCTION ferretdb_jsonb_to_sql_scalar(v jsonb) RETURNS varchar AS $$
-- BEGIN
-- END;
-- $$ LANGUAGE plpgsql IMMUTABLE;

-- CREATE OR REPLACE FUNCTION ferretdb_compare_sql_scalars(a anycompatiblenonarray, b anycompatiblenonarray) RETURNS smallint AS $$
-- BEGIN
--     PERFORM ferretdb_sql_scalar_type(a);
--     PERFORM ferretdb_sql_scalar_type(b);
--     RETURN CASE WHEN a = b THEN 0
--                 WHEN a > b THEN 1
--                 WHEN a < b THEN -1
--                 WHEN  a IS NULL AND b IS NULL THEN 0
--                 ELSE NULL
--            END;
-- END;
-- $$ LANGUAGE plpgsql IMMUTABLE;
