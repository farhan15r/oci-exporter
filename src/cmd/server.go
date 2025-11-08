package cmd

import (
	"net/http"
	"oci-exporter/src/config"
	"oci-exporter/src/handler"
	"oci-exporter/src/utils"
)

func StartServer() {
	config.InitConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handler.GETHome)

	mux.HandleFunc("GET /metrics", handler.GETMetrics)

	utils.Logger.Info("Starting Server on :" + config.Port)

	err := http.ListenAndServe(":"+config.Port, mux)
	if err != nil {
		utils.Logger.Error(err.Error())
		panic(err)
	}
}
