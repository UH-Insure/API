package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"


    "api/internal/auth"
    "api/internal/compiler"
    "api/internal/models"
    "api/internal/storage"
)


type API struct {
	Runner  *compiler.Runner
	Storage *storage.Storage
	APIKey  string
}

func NewAPI(r *compiler.Runner, s *storage.Storage, key string) *API {
	return &API{Runner: r, Storage: s, APIKey: key}
}

//
// ----- AUTH MIDDLEWARE -----
//

func (api *API) RequireKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if api.APIKey == "" {
			next(w, r)
			return
		}
		if r.Header.Get("X-API-Key") != api.APIKey {
			http.Error(w, "Unauthorized", 401)
			return
		}
		next(w, r)
	}
}

//
// ---- /run/{tool} ----
//

type RunRequest struct {
	Code string `json:"code"`
}

type RunResponse struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error,omitempty"`
	ID     string `json:"id"`
}

func (api *API) RunHandler(tool string) http.HandlerFunc {
	return api.RequireKey(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "POST only", 405)
			return
		}

		var body RunRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.Code == "" {
			http.Error(w, "missing code", 400)
			return
		}

		// --- Save input file ---
		f, err := api.Storage.SaveFile(tool+".input", []byte(body.Code))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// --- Run the tool ---
		stdout, stderr, execErr := api.Runner.Execute(tool, f.Path)

		errMsg := ""
		if execErr != nil {
			errMsg = execErr.Error()
		}

		api.Storage.SaveHistory(
			f.ID, f.Filename, f.Path,
			tool, stdout, stderr, errMsg,
		)

		json.NewEncoder(w).Encode(RunResponse{
			ID:     f.ID,
			Stdout: stdout,
			Stderr: stderr,
			Error:  errMsg,
		})
	})
}

//
// ---- File Upload ----
//

func (api *API) UploadHandler() http.HandlerFunc {
	return api.RequireKey(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer file.Close()

		data, _ := io.ReadAll(file)

		f, err := api.Storage.SaveFile(header.Filename, data)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id":       f.ID,
			"filename": f.Filename,
		})
	})
}

//
// ---- Get stored file ----
//

func (api *API) GetFile() http.HandlerFunc {
	return api.RequireKey(func(w http.ResponseWriter, r *http.Request) {
		id := filepath.Base(r.URL.Path)

		f, data, err := api.Storage.LoadFile(id)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+f.Filename)
		w.Write(data)
	})
}

//
// ---- History ----
//

func (api *API) HistoryHandler() http.HandlerFunc {
	return api.RequireKey(func(w http.ResponseWriter, r *http.Request) {
		h, err := api.Storage.ListHistory(50)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(h)
	})
}