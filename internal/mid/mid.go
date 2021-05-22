package mid

import "net/http"

// Middleware defines the standard framework agnostic middleware signature
type Middleware func(http.Handler) http.Handler
