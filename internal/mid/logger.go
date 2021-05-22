package mid

import (
	"log"
	"net/http"
	"time"
)

// Logger logs each request, once at the start and again at the end so we can track the latency of each request
func Logger(log *log.Logger) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := r.Context().Value(RequestValueKey).(*RequestValues)

			log.Printf("%s : started	: %s %s -> %s",
				v.TraceID,
				r.Method, r.URL.Path, r.RemoteAddr,
			)

			handler.ServeHTTP(w, r)

			log.Printf("%s : completed  : %s %s -> %s (%s)",
				v.TraceID,
				r.Method, r.URL.Path, r.RemoteAddr,
				time.Since(v.Now),
			)
		})
	}
}
