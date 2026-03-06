package middleware

import (
	"net/http"
)

// Auth is a middleware that validates basic auth
func Auth(password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, pass, ok := r.BasicAuth()
			if !ok || pass != password {
				w.Header().Set("WWW-Authenticate", `Basic realm="llm-eval"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
