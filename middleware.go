// Package middleware implements useful HTTP middleware for routers compatible
// with the standard library mux.
package middleware

import (
	"net/http"

	"github.com/holedaemon/slogx"
)

// Recoverer recovers from any panics in the request chain, logs the recover
// value and executes the given response handler.
func Recoverer(next http.Handler, fn http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				slogx.Error(r.Context(), "PANIC", "value", rvr)
				fn(w, r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RequestSize limits the size a request's body may be.
func RequestSize(bytes int64) func(http.Handler) http.Handler {
	f := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, bytes)
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
