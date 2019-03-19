package http

import "net/http"

func MaxBytesReader(bytes int64) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, bytes)
			handler.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
