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

func newPostgresqlMemoryUtilization() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "oci_exporter",
			Name:      "postgresql_memory_utilization",
			Help:      "Memory Utilization of OCI PostgreSQL Database.",
		},
		[]string{
			"resource_name",
			"resource_id",
			"compartment_id",
			"db_instance_id",
			"db_instance_role",
		},
	)
}

func GetPostgresqlMemoryUtilization(ctx context.Context, postgresqlMemoryUtilization *prometheus.GaugeVec) error {
	postgresqlMemoryUtilization.Reset()

	namespaceQuery := "oci_postgresql"
	query := "MemoryUtilization[1m].mean()"

	compartmentId := config.CompartmentId

	err := getPostgresqlMemoryUtilizationByCompartment(
		ctx,
		compartmentId,
		query,
		namespaceQuery,
		postgresqlMemoryUtilization,
	)
	if err != nil {
		return err
	}

	return nil
}

func getPostgresqlMemoryUtilizationByCompartment(
	ctx context.Context,
	compartmentId string,
	query string,
	namespaceQuery string,
	postgresqlMemoryUtilization *prometheus.GaugeVec,
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
		resourceName := metric.Dimensions["resourceName"]
		resourceId := metric.Dimensions["resourceId"]
		compartmentId := *metric.CompartmentId
		dbInstanceId := metric.Dimensions["dbInstanceId"]
		dbInstanceRole := metric.Dimensions["dbInstanceRole"]

		// set gauge value
		postgresqlMemoryUtilization.With(prometheus.Labels{
			"resource_name":    resourceName,
			"resource_id":      resourceId,
			"compartment_id":   compartmentId,
			"db_instance_id":   dbInstanceId,
			"db_instance_role": dbInstanceRole,
		}).Set(value)
	}

	return nil
}
