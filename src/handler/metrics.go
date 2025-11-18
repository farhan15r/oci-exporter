package handler

import (
	"net/http"
	"oci-exporter/src/services/oci"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	dbOracleExecuteCount prometheus.Collector
	dbOracleCurrLogon    prometheus.Collector

	dbExecuteCount prometheus.Collector
	dbCurrLogon    prometheus.Collector
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

	time.Sleep(1 * time.Second)

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

	dbOracleExecuteCount, err = oci.GetDbOracleExecuteCount(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	time.Sleep(1 * time.Second)

	dbOracleCurrLogon, err = oci.GetDbOracleCurrLogon(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	dbExecuteCount, err = oci.GetDbExecuteCount(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	dbCurrLogon, err = oci.GetDbCurrLogon(r.Context())
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
		prometheus.MustRegister(dbOracleExecuteCount)
		prometheus.MustRegister(dbOracleCurrLogon)
		prometheus.MustRegister(dbExecuteCount)
		prometheus.MustRegister(dbCurrLogon)
	})

	promhttp.Handler().ServeHTTP(w, r)
}
