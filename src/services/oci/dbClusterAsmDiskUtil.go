package oci

import (
	"context"
	"oci-exporter/src/config"
	"oci-exporter/src/utils"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/monitoring"

	"github.com/prometheus/client_golang/prometheus"
)

var dbClusterAsmDiskUtil = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "oci_exporter",
		Name:      "database_cluster_asm_disk_utilization",
		Help:      "ASM Disk Utilization of OCI Database Cluster.",
	},
	[]string{
		"reource_id",
		"compartment_id",
		"resource_name",
		"disk_group_name",
	},
)

func GetDbClusterAsmDiskUtil(ctx context.Context) (*prometheus.GaugeVec, error) {
	dbClusterAsmDiskUtil.Reset()

	namespaceQuery := "oci_database_cluster"
	query := "ASMDiskgroupUtilization[1m].max()"

	compartmentId := config.CompartmentId

	err := getDbClusterAsmDiskUtilByCompartment(
		ctx,
		compartmentId,
		query,
		namespaceQuery,
	)
	if err != nil {
		return nil, err
	}

	return dbClusterAsmDiskUtil, nil
}

func getDbClusterAsmDiskUtilByCompartment(
	ctx context.Context,
	compartmentId string,
	query string,
	namespaceQuery string,
) error {
	minutes := config.TimeRangeMinute

	end := time.Now().UTC()
	start := end.Add(-time.Duration(minutes) * time.Minute)

	sdkStart := common.SDKTime{Time: start}
	sdkEnd := common.SDKTime{Time: end}

	req := monitoring.SummarizeMetricsDataRequest{
		CompartmentId:          &compartmentId,
		CompartmentIdInSubtree: &config.CompartmentIdInSubtree,
		SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
			Query:     &query,
			StartTime: &sdkStart,
			EndTime:   &sdkEnd,
			Namespace: &namespaceQuery,
		},
	}

	client, err := config.NewOciClient()
	if err != nil {
		utils.Logger.Error("failed to create OCI client", "error", err.Error())
		return err
	}

	resp, err := client.SummarizeMetricsData(ctx, req)
	if err != nil {
		utils.Logger.Error("SummarizeMetricsData failed", "error", err.Error())
		return err
	}

	for _, metric := range resp.Items {
		// get last data point
		if len(metric.AggregatedDatapoints) == 0 {
			continue
		}
		lastPoint := metric.AggregatedDatapoints[len(metric.AggregatedDatapoints)-1]
		value := *lastPoint.Value

		// extract dimension values
		resourceId := metric.Dimensions["resourceId"]
		diskgroupName := metric.Dimensions["diskgroupName"]
		resourceName := metric.Dimensions["resourceName"]
		compartmentId := *metric.CompartmentId

		// set gauge value
		dbClusterAsmDiskUtil.With(prometheus.Labels{
			"reource_id":      resourceId,
			"compartment_id":  compartmentId,
			"resource_name":   resourceName,
			"disk_group_name": diskgroupName,
		}).Set(value)
	}

	return nil
}
