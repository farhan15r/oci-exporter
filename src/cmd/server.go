package cmd

import (
	"net/http"
	"oci-exporter/src/handler"
	"oci-exporter/src/utils"
	"os"
)

func StartServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handler.GETHome)

	mux.HandleFunc("GET /metrics", handler.GETMetrics)

	utils.Logger.Info("Starting Server on :" + port)

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		utils.Logger.Error(err.Error())
		panic(err)
	}
}
