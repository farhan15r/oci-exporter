package handler

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"oci-exporter/src/services/oci"
)

var (
	registerOnce sync.Once

	fastconnectBgpSession       prometheus.Collector
	fastconnectBytesReceivedSum prometheus.Collector
	fastconnectBytesSentSum     prometheus.Collector

	vpnBgpSession       prometheus.Collector
	vpnIpSecState       prometheus.Collector
	vpnBytesReceivedSum prometheus.Collector
	vpnBytesSentSum     prometheus.Collector

	dbClusterAsmDiskUtil prometheus.Collector
	dbClusterNodeStatus  prometheus.Collector
)

func GETMetrics(w http.ResponseWriter, r *http.Request) {
	err := error(nil)

	fastconnectBgpSession, err = oci.GetFastconnectBGPSessionState(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	fastconnectBytesReceivedSum, err = oci.GetFastconnectBytesReceivedSum(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	fastconnectBytesSentSum, err = oci.GetFastconnectBytesSentSum(r.Context())
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

	vpnIpSecState, err = oci.GetVpnIpSecState(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	vpnBytesReceivedSum, err = oci.GetVpnBytesReceivedSum(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	vpnBytesSentSum, err = oci.GetVpnBytesSentSum(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	dbClusterAsmDiskUtil, err = oci.GetDbClusterAsmDiskUtil(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	dbClusterNodeStatus, err = oci.GetDbClusterNodeStatus(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	registerOnce.Do(func() {
		prometheus.MustRegister(fastconnectBgpSession)
		prometheus.MustRegister(fastconnectBytesReceivedSum)
		prometheus.MustRegister(fastconnectBytesSentSum)
		prometheus.MustRegister(vpnBgpSession)
		prometheus.MustRegister(vpnIpSecState)
		prometheus.MustRegister(vpnBytesReceivedSum)
		prometheus.MustRegister(vpnBytesSentSum)
		prometheus.MustRegister(dbClusterAsmDiskUtil)
		prometheus.MustRegister(dbClusterNodeStatus)
	})

	promhttp.Handler().ServeHTTP(w, r)
}
