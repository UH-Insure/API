#!/bin/sh
set -e

DB_PATH="/data/sawapi.db"
SCHEMA_PATH="/data/schema.sql"

echo "[init-db] Checking database..."

if [ ! -f "$DB_PATH" ]; then
    echo "[init-db] No database found. Creating new SQLite DB at $DB_PATH ..."
    
    if [ ! -f "$SCHEMA_PATH" ]; then
        echo "[init-db] ERROR: schema.sql not found at $SCHEMA_PATH"
        exit 1
    fi

    sqlite3 "$DB_PATH" < "$SCHEMA_PATH"
    echo "[init-db] Database created successfully."
else
    echo "[init-db] Database already exists. Skipping initialization."
fi