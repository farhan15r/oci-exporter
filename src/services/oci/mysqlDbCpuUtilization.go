package oci

import (
	"context"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/monitoring"
	"github.com/prometheus/client_golang/prometheus"

	"oci-exporter/src/config"
	"oci-exporter/src/utils"
)

func newMysqlDbCpuUtilization() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "oci_exporter",
			Name:      "mysql_db_cpu_utilization",
			Help:      "CPU Utilization of OCI MySQL Database.",
		},
		[]string{
			"compartment_id",
			"resource_name",
			"resource_id",
			"resource_type",
			"heat_wave_node",
		},
	)
}

func GetMysqlDbCpuUtilization(ctx context.Context, mysqlDbCpuUtilization *prometheus.GaugeVec) error {
	mysqlDbCpuUtilization.Reset()

	namespaceQuery := "oci_mysql_database"
	query := "CPUUtilization[1m].mean()"

	compartmentId := config.CompartmentId

	err := getMysqlDbCpuUtilizationByCompartment(
		ctx,
		compartmentId,
		query,
		namespaceQuery,
		mysqlDbCpuUtilization,
	)
	if err != nil {
		return err
	}

	return nil
}

func getMysqlDbCpuUtilizationByCompartment(
	ctx context.Context,
	compartmentId string,
	query string,
	namespaceQuery string,
	mysqlDbCpuUtilization *prometheus.GaugeVec,
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
		compartmentId := *metric.CompartmentId
		resourceName := metric.Dimensions["resourceName"]
		resourceId := metric.Dimensions["resourceId"]
		resourceType := metric.Dimensions["resourceType"]
		heatWaveNode := metric.Dimensions["heatWaveNode"]

		// set gauge value
		mysqlDbCpuUtilization.With(prometheus.Labels{
			"resource_name":  resourceName,
			"resource_id":    resourceId,
			"compartment_id": compartmentId,
			"resource_type":  resourceType,
			"heat_wave_node": heatWaveNode,
		}).Set(value)
	}

	return nil
}
