package helpers

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TraceHandler(method, path string, h http.HandlerFunc) http.Handler {
	return otelhttp.NewHandler(h, fmt.Sprintf("%s %s", method, path))
}
