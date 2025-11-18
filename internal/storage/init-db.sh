#!/bin/sh
set -e

DB_PATH="/data/sawapi.db"
SCHEMA_PATH="/data/schema.sql"

echo "[init-db] Starting initialization..."

# Debug print
echo "[init-db] Checking if schema exists at: $SCHEMA_PATH"
if [ ! -f "$SCHEMA_PATH" ]; then
    echo "[init-db] ERROR: schema.sql not found at $SCHEMA_PATH"
    ls -l /data || true
    exit 1
fi

# Check DB existence
if [ ! -f "$DB_PATH" ]; then
    echo "[init-db] No database found. Creating new SQLite DB at $DB_PATH ..."
    sqlite3 "$DB_PATH" < "$SCHEMA_PATH"
    echo "[init-db] Database created successfully."
else
    echo "[init-db] Database already exists. Skipping creation."
fi

echo "[init-db] Launching API server..."
exec sawapi