package tensile

import (
	"net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler
