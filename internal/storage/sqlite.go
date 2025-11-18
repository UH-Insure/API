package storage

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func InitSQLite(path string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS history (
            id TEXT PRIMARY KEY,
            filename TEXT,
            path TEXT,
            tool TEXT,
            stdout TEXT,
            stderr TEXT,
            error TEXT,
            created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `)
    return db, err
}