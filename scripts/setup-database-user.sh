#!/bin/bash

# Database User Setup Script for godev-kit
# This script sets up a PostgreSQL user with full database access

set -e

# Configuration
DB_NAME="godevkit"
DB_USER="godevkit"
DB_PASSWORD="your_secure_password"  # Change this to a secure password
PG_HOST="localhost"
PG_PORT="5433"  # Based on your postgresql.conf
PG_SUPERUSER="postgres"

echo "🚀 Setting up database user for godev-kit..."

# Check if PostgreSQL is running
if ! pg_isready -h $PG_HOST -p $PG_PORT -U $PG_SUPERUSER > /dev/null 2>&1; then
    echo "❌ PostgreSQL is not running on $PG_HOST:$PG_PORT"
    echo "Please start PostgreSQL first:"
    echo "  brew services start postgresql"
    echo "  or"
    echo "  pg_ctl -D /usr/local/var/postgres start"
    exit 1
fi

echo "✅ PostgreSQL is running"

# Create the SQL script with the actual password
cat > /tmp/setup_db_user.sql << EOF
-- Database User Setup Script for godev-kit
-- Run this script as a PostgreSQL superuser (postgres)

-- 1. Create the database if it doesn't exist
CREATE DATABASE $DB_NAME;

-- 2. Create user if it doesn't exist
DO \$\$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '$DB_USER') THEN
        CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';
    END IF;
END
\$\$;

-- 3. Grant all privileges on the database
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;

-- 4. Connect to the database and grant schema privileges
\c $DB_NAME;

-- 5. Grant usage on public schema
GRANT USAGE ON SCHEMA public TO $DB_USER;

-- 6. Grant all privileges on all existing tables
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO $DB_USER;

-- 7. Grant all privileges on all existing sequences
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO $DB_USER;

-- 8. Grant all privileges on all existing functions
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO $DB_USER;

-- 9. Grant create privileges on schema
GRANT CREATE ON SCHEMA public TO $DB_USER;

-- 10. Set default privileges for future objects
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $DB_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO $DB_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO $DB_USER;

-- 11. Verify the setup
\du $DB_USER;
EOF

echo "📝 Executing database setup script..."

# Execute the SQL script
psql -h $PG_HOST -p $PG_PORT -U $PG_SUPERUSER -d postgres -f /tmp/setup_db_user.sql

# Clean up
rm -f /tmp/setup_db_user.sql

echo "✅ Database user setup completed!"
echo ""
echo "📋 Summary:"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo "  Host: $PG_HOST"
echo "  Port: $PG_PORT"
echo ""
echo "🔗 Connection string:"
echo "  postgres://$DB_USER:$DB_PASSWORD@$PG_HOST:$PG_PORT/$DB_NAME?sslmode=disable"
echo ""
echo "⚠️  Remember to update your config files with the correct password!" 