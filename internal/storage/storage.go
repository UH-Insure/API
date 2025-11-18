package storage

import (
    "database/sql"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    _ "github.com/mattn/go-sqlite3"
    "github.com/google/uuid"
)

type Storage struct {
	DB *sql.DB
}

type StoredFile struct {
	ID       string
	Path     string
	Filename string
}

func NewStorage(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	s := &Storage{DB: db}
	err = s.init()
	return s, err
}

func (s *Storage) init() error {
	_, err := s.DB.Exec(`
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
	return err
}

//
// ---- FILE STORAGE ----
//

func (s *Storage) SaveFile(originalName string, data []byte) (*StoredFile, error) {
	id := uuid.New().String()

	dir := "/work/files"
	os.MkdirAll(dir, 0755)

	dst := filepath.Join(dir, id+"_"+originalName)
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return nil, err
	}

	return &StoredFile{
		ID:       id,
		Path:     dst,
		Filename: originalName,
	}, nil
}

func (s *Storage) LoadFile(id string) (*StoredFile, []byte, error) {
	dir := "/work/files"

	matches, _ := filepath.Glob(filepath.Join(dir, id+"_*"))
	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("file not found")
	}

	path := matches[0]
	base := filepath.Base(path)
	_, filename, _ := strings.Cut(base, "_")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	return &StoredFile{
		ID:       id,
		Path:     path,
		Filename: filename,
	}, data, nil
}

//
// ---- HISTORY LOG ----
//

func (s *Storage) SaveHistory(id, filename, path, tool, stdout, stderr, errStr string) error {
	_, err := s.DB.Exec(`
		INSERT INTO history (id, filename, path, tool, stdout, stderr, error)
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`, id, filename, path, tool, stdout, stderr, errStr)
	return err
}

func (s *Storage) ListHistory(limit int) ([]map[string]any, error) {
	rows, err := s.DB.Query(`
		SELECT id, filename, tool, created
		FROM history
		ORDER BY created DESC
		LIMIT ?;
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []map[string]any

	for rows.Next() {
		var id, filename, tool string
		var created string
		rows.Scan(&id, &filename, &tool, &created)

		res = append(res, map[string]any{
			"id":       id,
			"filename": filename,
			"tool":     tool,
			"created":  created,
		})
	}

	return res, nil
}