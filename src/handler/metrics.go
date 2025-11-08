package handler

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"oci-exporter/src/services/oci"
)

var (
	registerOnce             sync.Once
	fastconnectBgpSession    prometheus.Collector
	fastconnectBytesReceived prometheus.Collector
	vpnBgpSession            prometheus.Collector
)

func GETMetrics(w http.ResponseWriter, r *http.Request) {
	err := error(nil)

	fastconnectBgpSession, err = oci.GetFastconnectBGPSessionState(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	fastconnectBytesReceived, err = oci.GetFastconnectBytesReceived(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	vpnBgpSession, err = oci.GetVpnBGPSessionState(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	registerOnce.Do(func() {
		prometheus.MustRegister(fastconnectBgpSession)
		prometheus.MustRegister(fastconnectBytesReceived)
		prometheus.MustRegister(vpnBgpSession)
	})

	promhttp.Handler().ServeHTTP(w, r)
}
