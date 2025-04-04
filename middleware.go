// Package middleware implements useful HTTP middleware for routers compatible
// with the standard library mux.
package middleware

import (
	"net/http"

	"github.com/holedaemon/slogx"
)

// Recoverer recovers from any panics in the request chain, logs the recover
// value and executes the given response handler.
func Recoverer(recoverFunc http.HandlerFunc) func(next http.Handler) http.Handler {
	f := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					slogx.Error(r.Context(), "PANIC", "value", rvr)
					recoverFunc(w, r)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return f
}

// CORS allows cross-site requests for the given origin.
func CORS(origin string) func(http.Handler) http.Handler {
	f := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}

	return f
}
