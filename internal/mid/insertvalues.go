package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ctxKey int

var RequestValueKey ctxKey = 0

type RequestValues struct {
	TraceID string
	Now     time.Time
}

// InsertValues places RequestValues in the context for each request so we can access the contents in handlers/resolvers
func InsertValues() Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := RequestValues{
				TraceID: uuid.New().String(),
				Now:     time.Now(),
			}
			ctx := context.WithValue(r.Context(), RequestValueKey, &v)
			handler.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
