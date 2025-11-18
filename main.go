package main

import (
	"log"
	"net/http"
	"os"

    "api/api"
    "api/internal/compiler"
    "api/internal/storage"
)

func main() {
	// --- Load API key ---
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "dev-key"
	}

	// --- Init storage ---
	store, err := storage.NewStorage("/work/history.db")
	if err != nil {
		log.Fatal(err)
	}

	// --- Runner ---
	runner := compiler.NewRunner()

	// --- API ---
	app := api.NewAPI(runner, store, apiKey)

	// ---- ROUTES ----
	http.Handle("/run/cryptol", app.RunHandler("cryptol"))
	http.Handle("/run/saw", app.RunHandler("saw"))
	http.Handle("/run/c", app.RunHandler("c"))

	http.Handle("/store/upload", app.UploadHandler())
	http.Handle("/store/get/", http.StripPrefix("/store/get/", app.GetFile()))
	http.Handle("/history", app.HistoryHandler())

	log.Println("API listening on :8443")
	log.Fatal(http.ListenAndServe(":8443", nil))
}