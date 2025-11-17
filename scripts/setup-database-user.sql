-- Database User Setup Script for godev-kit
-- Run this script as a PostgreSQL superuser (postgres)

-- 1. Create the database if it doesn't exist
CREATE DATABASE godevkit;

-- 2. Create user if it doesn't exist (replace 'your_secure_password' with actual password)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'godevkit') THEN
        CREATE USER godevkit WITH PASSWORD 'your_secure_password';
    END IF;
END
$$;

-- 3. Grant all privileges on the database
GRANT ALL PRIVILEGES ON DATABASE godevkit TO godevkit;

-- 4. Connect to the godevkit database and grant schema privileges
\c godevkit;

-- 5. Grant usage on public schema
GRANT USAGE ON SCHEMA public TO godevkit;

-- 6. Grant all privileges on all existing tables
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO godevkit;

-- 7. Grant all privileges on all existing sequences
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO godevkit;

-- 8. Grant all privileges on all existing functions
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO godevkit;

-- 9. Grant create privileges on schema
GRANT CREATE ON SCHEMA public TO godevkit;

-- 10. Set default privileges for future objects
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO godevkit;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO godevkit;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO godevkit;

-- 11. Verify the setup
\du godevkit; 