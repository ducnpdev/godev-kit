-- Grant All Schemas Access Script
-- This script grants full access to ALL schemas in the database
-- Run this as a PostgreSQL superuser

-- Configuration
DO $$
DECLARE
    schema_record RECORD;
    user_name TEXT := 'godevkit';  -- Change this to your target user
BEGIN
    -- Grant usage on all schemas
    FOR schema_record IN 
        SELECT schema_name 
        FROM information_schema.schemata 
        WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast')
    LOOP
        EXECUTE format('GRANT USAGE ON SCHEMA %I TO %I', schema_record.schema_name, user_name);
        RAISE NOTICE 'Granted USAGE on schema: %', schema_record.schema_name;
    END LOOP;

    -- Grant all privileges on all tables in all schemas
    FOR schema_record IN 
        SELECT schema_name 
        FROM information_schema.schemata 
        WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast')
    LOOP
        EXECUTE format('GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA %I TO %I', schema_record.schema_name, user_name);
        RAISE NOTICE 'Granted ALL PRIVILEGES on ALL TABLES in schema: %', schema_record.schema_name;
    END LOOP;

    -- Grant all privileges on all sequences in all schemas
    FOR schema_record IN 
        SELECT schema_name 
        FROM information_schema.schemata 
        WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast')
    LOOP
        EXECUTE format('GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA %I TO %I', schema_record.schema_name, user_name);
        RAISE NOTICE 'Granted ALL PRIVILEGES on ALL SEQUENCES in schema: %', schema_record.schema_name;
    END LOOP;

    -- Grant all privileges on all functions in all schemas
    FOR schema_record IN 
        SELECT schema_name 
        FROM information_schema.schemata 
        WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast')
    LOOP
        EXECUTE format('GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA %I TO %I', schema_record.schema_name, user_name);
        RAISE NOTICE 'Granted ALL PRIVILEGES on ALL FUNCTIONS in schema: %', schema_record.schema_name;
    END LOOP;

    -- Grant create privileges on all schemas (if user should be able to create objects)
    FOR schema_record IN 
        SELECT schema_name 
        FROM information_schema.schemata 
        WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast')
    LOOP
        EXECUTE format('GRANT CREATE ON SCHEMA %I TO %I', schema_record.schema_name, user_name);
        RAISE NOTICE 'Granted CREATE on schema: %', schema_record.schema_name;
    END LOOP;

    -- Set default privileges for future objects in all schemas
    FOR schema_record IN 
        SELECT schema_name 
        FROM information_schema.schemata 
        WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast')
    LOOP
        EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL ON TABLES TO %I', schema_record.schema_name, user_name);
        EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL ON SEQUENCES TO %I', schema_record.schema_name, user_name);
        EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL ON FUNCTIONS TO %I', schema_record.schema_name, user_name);
        RAISE NOTICE 'Set DEFAULT PRIVILEGES for future objects in schema: %', schema_record.schema_name;
    END LOOP;

END $$;

-- Show all schemas and their privileges
SELECT 
    schema_name,
    'USAGE' as privilege_type
FROM information_schema.schemata 
WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast')
UNION ALL
SELECT 
    schema_name,
    'CREATE' as privilege_type
FROM information_schema.schemata 
WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast');

-- Verify user privileges
\du godevkit; 