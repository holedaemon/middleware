package middleware

import (
	"log/slog"
	"net/http"

	"github.com/holedaemon/slogx"
)

// Logger adds the given *[slog.Logger] to a request's [context.Context].
func Logger(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := slogx.WithLogger(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// RequestLogger logs any incoming requests to the debug level if a
// *[slog.Logger] is present in the request's [context.Context].
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slogx.Debug(r.Context(), "http request",
			"method", r.Method,
			"url", r.RequestURI,
			"proto", r.Proto,
		)

		next.ServeHTTP(w, r)
	})
}
