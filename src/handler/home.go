package handler

import (
	"encoding/json"
	"net/http"
)

func GETHome(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"status":           "success",
		"message":          "Welcome to oci-exporter",
		"metrics_endpoint": "/metrics",
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error generating response: " + err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
