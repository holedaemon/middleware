package middleware

import (
	"context"
	"net/http"

	"github.com/holedaemon/slogx"
	"github.com/rs/xid"
)

// RequestIDHeader is the key used for request IDs.
const RequestIDHeader = "X-Request-ID"

type requestIDKey struct{}

// RequestID adds a unique identifier to each request.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var id xid.ID
		rid := r.Header.Get(RequestIDHeader)
		if rid == "" {
			id = xid.New()
			rid = id.String()
		} else {
			var err error
			id, err = xid.FromString(rid)
			if err != nil {
				old := rid
				id = xid.New()
				rid = id.String()

				slogx.Debug(ctx, "replacing request ID", "old", old, "new", rid)
			}
		}

		w.Header().Set(RequestIDHeader, rid)
		ctx = context.WithValue(ctx, requestIDKey{}, id)
		ctx = slogx.With(ctx, "request_id", rid)

		r = r.WithContext(ctx)
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
