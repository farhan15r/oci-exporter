package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GETMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
