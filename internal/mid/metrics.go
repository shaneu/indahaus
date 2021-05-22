package mid

import (
	"expvar"
	"net/http"
	"runtime"
)

// m contains the global program counters for the application
var m = struct {
	gr  *expvar.Int
	req *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
}

// Metrics collects information about the number of requests and goroutines for viewing at host:port/debug/vars
func Metrics() Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// increment the request counter
			m.req.Add(1)

			// every 100 requests get the number of goroutines, calling runtime.NumGoroutine isn't
			// free so we avoid doing it on each request
			if m.req.Value()%100 == 0 {
				m.gr.Set(int64(runtime.NumGoroutine()))
			}

			handler.ServeHTTP(w, r)
		})
	}
}
