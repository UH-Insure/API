-- SQLite schema for SAW/Cryptol API

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