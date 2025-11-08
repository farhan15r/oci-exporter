package config

import (
	"oci-exporter/src/utils"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/monitoring"
)

func NewOciClient() (monitoring.MonitoringClient, error) {
	confProvider := common.CustomProfileConfigProvider("./oci.config", "DEFAULT")
	client, err := monitoring.NewMonitoringClientWithConfigurationProvider(confProvider)
	if err != nil {
		utils.Logger.Error("failed creating monitoring client", "error", err.Error())
		return monitoring.MonitoringClient{}, err
	}

	return client, nil
}
