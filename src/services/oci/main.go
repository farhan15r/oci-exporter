package oci

import (
	"context"
	"oci-exporter/src/config"
	"oci-exporter/src/utils"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	dbClusterAsmDiskUtil        = newDbClusterAsmDiskUtil()
	dbClusterNodeStatus         = newDbClusterNodeStatus()
	dbCurrLogon                 = newDbCurrLogon()
	dbExecuteCount              = newDbExecuteCount()
	dbOracleCurrLogon           = newDbOracleCurrLogon()
	dbOracleExecuteCount        = newDbOracleExecuteCount()
	fastconnectBgpSession       = newFastconnectBgpSession()
	fastconnectBytesReceivedSum = newFastconnectBytesReceivedSum()
	fastconnectBytesSentSum     = newFastconnectBytesSentSum()
	vpnBgpSession               = newVpnBgpSession()
	vpnBytesReceivedSum         = newVpnBytesReceivedSum()
	vpnBytesSentSum             = newVpnBytesSentSum()
	vpnIpSecState               = newVpnIpSecState()
)

func InitAndRegister() {
	var err error

	// dbClusterAsmDiskUtil
	err = GetDbClusterAsmDiskUtil(context.Background(), dbClusterAsmDiskUtil)
	if err != nil {
		utils.Logger.Error("GetDbClusterAsmDiskUtil init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(dbClusterAsmDiskUtil)
	}

	// dbClusterNodeStatus
	err = GetDbClusterNodeStatus(context.Background(), dbClusterNodeStatus)
	if err != nil {
		utils.Logger.Error("GetDbClusterNodeStatus init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(dbClusterNodeStatus)
	}

	// dbCurrLogon
	err = GetDbCurrLogon(context.Background(), dbCurrLogon)
	if err != nil {
		utils.Logger.Error("GetDbCurrLogon init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(dbCurrLogon)
	}

	// dbExecuteCount
	err = GetDbExecuteCount(context.Background(), dbExecuteCount)
	if err != nil {
		utils.Logger.Error("GetDbExecuteCount init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(dbExecuteCount)
	}

	// dbOracleCurrLogon
	err = GetDbOracleCurrLogon(context.Background(), dbOracleCurrLogon)
	if err != nil {
		utils.Logger.Error("GetDbOracleCurrLogon init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(dbOracleCurrLogon)
	}

	// dbOracleExecuteCount
	err = GetDbOracleExecuteCount(context.Background(), dbOracleExecuteCount)
	if err != nil {
		utils.Logger.Error("GetDbOracleExecuteCount init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(dbOracleExecuteCount)
	}

	// fastconnectBgpSession
	err = GetFastconnectBGPSessionState(context.Background(), fastconnectBgpSession)
	if err != nil {
		utils.Logger.Error("GetFastconnectBGPSessionState init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(fastconnectBgpSession)
	}

	// fastconnectBytesReceivedSum
	err = GetFastconnectBytesReceivedSum(context.Background(), fastconnectBytesReceivedSum)
	if err != nil {
		utils.Logger.Error("GetFastconnectBytesReceivedSum init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(fastconnectBytesReceivedSum)
	}

	// fastconnectBytesSentSum
	err = GetFastconnectBytesSentSum(context.Background(), fastconnectBytesSentSum)
	if err != nil {
		utils.Logger.Error("GetFastconnectBytesSentSum init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(fastconnectBytesSentSum)
	}

	// vpnBgpSession
	err = GetVpnBGPSessionState(context.Background(), vpnBgpSession)
	if err != nil {
		utils.Logger.Error("GetVpnBGPSessionState init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(vpnBgpSession)
	}

	// vpnBytesReceivedSum
	err = GetVpnBytesReceivedSum(context.Background(), vpnBytesReceivedSum)
	if err != nil {
		utils.Logger.Error("GetVpnBytesReceivedSum init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(vpnBytesReceivedSum)
	}

	// vpnBytesSentSum
	err = GetVpnBytesSentSum(context.Background(), vpnBytesSentSum)
	if err != nil {
		utils.Logger.Error("GetVpnBytesSentSum init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(vpnBytesSentSum)
	}

	// vpnIpSecState
	err = GetVpnIpSecState(context.Background(), vpnIpSecState)
	if err != nil {
		utils.Logger.Error("GetVpnIpSecState init failed", "error", err.Error())
	} else {
		prometheus.MustRegister(vpnIpSecState)
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
			var err error

			err = GetDbClusterAsmDiskUtil(ctx, dbClusterAsmDiskUtil)
			if err != nil {
				utils.Logger.Error("refresh GetDbClusterAsmDiskUtil failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetDbClusterNodeStatus(ctx, dbClusterNodeStatus)
			if err != nil {
				utils.Logger.Error("refresh GetDbClusterNodeStatus failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetDbCurrLogon(ctx, dbCurrLogon)
			if err != nil {
				utils.Logger.Error("refresh GetDbCurrLogon failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetDbExecuteCount(ctx, dbExecuteCount)
			if err != nil {
				utils.Logger.Error("refresh GetDbExecuteCount failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetDbOracleCurrLogon(ctx, dbOracleCurrLogon)
			if err != nil {
				utils.Logger.Error("refresh GetDbOracleCurrLogon failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetDbOracleExecuteCount(ctx, dbOracleExecuteCount)
			if err != nil {
				utils.Logger.Error("refresh GetDbOracleExecuteCount failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetFastconnectBGPSessionState(ctx, fastconnectBgpSession)
			if err != nil {
				utils.Logger.Error("refresh GetFastconnectBGPSessionState failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetFastconnectBytesReceivedSum(ctx, fastconnectBytesReceivedSum)
			if err != nil {
				utils.Logger.Error("refresh GetFastconnectBytesReceivedSum failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetFastconnectBytesSentSum(ctx, fastconnectBytesSentSum)
			if err != nil {
				utils.Logger.Error("refresh GetFastconnectBytesSentSum failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetVpnBGPSessionState(ctx, vpnBgpSession)
			if err != nil {
				utils.Logger.Error("refresh GetVpnBGPSessionState failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetVpnBytesReceivedSum(ctx, vpnBytesReceivedSum)
			if err != nil {
				utils.Logger.Error("refresh GetVpnBytesReceivedSum failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetVpnBytesSentSum(ctx, vpnBytesSentSum)
			if err != nil {
				utils.Logger.Error("refresh GetVpnBytesSentSum failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)

			err = GetVpnIpSecState(ctx, vpnIpSecState)
			if err != nil {
				utils.Logger.Error("refresh GetVpnIpSecState failed", "error", err.Error())
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}
