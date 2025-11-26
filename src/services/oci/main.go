package oci

import (
	"context"
	"oci-exporter/src/config"
	"oci-exporter/src/utils"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	dbClusterAsmDiskUtil = newDbClusterAsmDiskUtil()
)

func InitAndRegister() {
	var err error

	err = GetDbClusterAsmDiskUtil(context.Background(), dbClusterAsmDiskUtil)
	if err != nil {
		utils.Logger.Error("GetDbClusterAsmDiskUtil init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(dbClusterAsmDiskUtil)
	}
}

func StartBackgroundUpdater(ctx context.Context) {
	interval := time.Duration(config.RefreshIntervalSeconds) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	utils.Logger.Info("Starting background metrics updater", "interval_seconds", config.RefreshIntervalSeconds)

	for {
		select {
		case <-ctx.Done():
			utils.Logger.Info("Background updater stopping: context done")
			return
		case <-ticker.C:
			err := GetDbClusterAsmDiskUtil(ctx, dbClusterAsmDiskUtil)
			if err != nil {
				utils.Logger.Error("refresh GetDbClusterAsmDiskUtil failed", "error", err.Error())
			}
		}
	}
}
