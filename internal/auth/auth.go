package auth

import "net/http"

func RequireKeyMiddleware(apiKey string, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if apiKey != "" && r.Header.Get("X-API-Key") != apiKey {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}