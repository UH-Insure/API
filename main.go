package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type RunRequest struct {
	FilePath string   `json:"file_path"`
	Code     string   `json:"code,omitempty"` // NEW: allows inline source code
	Args     []string `json:"args,omitempty"`
}

type RunResponse struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error,omitempty"`
}

func runCommand(tool string, args []string) (string, string, error) {
	cmd := exec.Command(tool, args...)
	cmd.Dir = "/work"
	cmd.Env = append(os.Environ(), "PATH=/usr/local/bin:/usr/bin:/bin")

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Start(); err != nil {
		return "", "", err
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-time.After(30 * time.Second):
		_ = cmd.Process.Kill()
		return outBuf.String(), errBuf.String(), os.ErrDeadlineExceeded
	case err := <-done:
		return outBuf.String(), errBuf.String(), err
	}
}

// handleRun wraps command execution for /run/{tool} endpoints
func handleRun(tool string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Use POST", http.StatusMethodNotAllowed)
			return
		}

		// Optional API key check
		apiKey := os.Getenv("API_KEY")
		if apiKey != "" && r.Header.Get("X-API-Key") != apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var req RunRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Decide input file path
		path := req.FilePath
		if req.Code != "" {
			// Create a temp file under /work
			tmpDir := "/work"
			if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
				os.MkdirAll(tmpDir, 0755)
			}
			tmpFile, err := os.CreateTemp(tmpDir, tool+"-*.input")
			if err != nil {
				http.Error(w, "Failed to create temp file: "+err.Error(), 500)
				return
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(req.Code); err != nil {
				http.Error(w, "Failed to write temp file: "+err.Error(), 500)
				return
			}
			tmpFile.Close()

			path = tmpFile.Name()
		}

		if path == "" {
			http.Error(w, "Missing file_path or code", http.StatusBadRequest)
			return
		}

		args := append([]string{filepath.Base(path)}, req.Args...)
		stdout, stderr, err := runCommand(tool, args)

		resp := RunResponse{Stdout: stdout, Stderr: stderr}
		if err != nil {
			resp.Error = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func main() {
	http.HandleFunc("/run/cryptol", handleRun("cryptol"))
	http.HandleFunc("/run/saw", handleRun("saw"))

	port := ":8443"
	useTLS := false // toggle TLS

	if useTLS {
		cert := "server.crt"
		key := "server.key"
		if _, err := os.Stat(cert); os.IsNotExist(err) {
			log.Println("No TLS certs found — generating self-signed certs...")
			exec.Command("openssl", "req", "-x509", "-newkey", "rsa:2048",
				"-nodes", "-keyout", key, "-out", cert,
				"-subj", "/CN=localhost", "-days", "365").Run()
		}
		log.Println("SAW/Cryptol API running on https://0.0.0.0" + port)
		log.Fatal(http.ListenAndServeTLS(port, cert, key, nil))
	} else {
		log.Println("SAW/Cryptol API running on http://0.0.0.0" + port)
		log.Fatal(http.ListenAndServe(port, nil))
	}
}