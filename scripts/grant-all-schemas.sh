#!/bin/bash

# Grant All Schemas Access Script
# This script grants full access to ALL schemas in the database

set -e

# Configuration
DB_NAME="godevkit"
DB_USER="godevkit"  # Change this to your target user
PG_HOST="localhost"
PG_PORT="5433"
PG_SUPERUSER="postgres"

echo "🚀 Granting access to ALL schemas for user: $DB_USER"

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

# Check if user exists
if ! psql -h $PG_HOST -p $PG_PORT -U $PG_SUPERUSER -d $DB_NAME -t -c "SELECT 1 FROM pg_roles WHERE rolname='$DB_USER'" | grep -q 1; then
    echo "❌ User '$DB_USER' does not exist"
    echo "Please create the user first or run: ./scripts/setup-database-user.sh"
    exit 1
fi

echo "✅ User '$DB_USER' exists"

# Execute the schema grant script
echo "📝 Granting privileges on all schemas..."

psql -h $PG_HOST -p $PG_PORT -U $PG_SUPERUSER -d $DB_NAME -f scripts/grant-all-schemas.sql

echo "✅ All schema privileges granted successfully!"
echo ""
echo "📋 Summary of granted privileges:"
echo "  ✅ USAGE on all schemas"
echo "  ✅ ALL PRIVILEGES on all tables in all schemas"
echo "  ✅ ALL PRIVILEGES on all sequences in all schemas"
echo "  ✅ ALL PRIVILEGES on all functions in all schemas"
echo "  ✅ CREATE on all schemas"
echo "  ✅ DEFAULT PRIVILEGES for future objects"
echo ""
echo "🔍 To verify, you can run:"
echo "  psql -h $PG_HOST -p $PG_PORT -U $DB_USER -d $DB_NAME -c '\du'" 