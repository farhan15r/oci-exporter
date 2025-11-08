package handler

import (
	"net/http"
	"oci-exporter/src/services/oci"
)

func GETMetrics(w http.ResponseWriter, r *http.Request) {
	// instanceCPUUtilization, err := oci.GetInstanceCPUUtilization(r.Context())
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte("Error fetching metrics: " + err.Error()))
	// 	return
	// }

	fastConnBgpSess, err := oci.GetFastconnectBGPSessionState(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	fastBytesReceived, err := oci.GetFastconnectBytesReceived(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	vpnBgpSess, err := oci.GetVpnBGPSessionState(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching metrics: " + err.Error()))
		return
	}

	res := ""
	// res += instanceCPUUtilization
	res += fastConnBgpSess
	res += fastBytesReceived
	res += vpnBgpSess

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(res))
}
