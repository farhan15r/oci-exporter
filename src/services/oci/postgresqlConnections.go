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

func newPostgresqlConnections() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "oci_exporter",
			Name:      "postgresql_connections",
			Help:      "Connections of OCI PostgreSQL Database.",
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

func GetPostgresqlConnections(ctx context.Context, postgresqlConnections *prometheus.GaugeVec) error {
	postgresqlConnections.Reset()

	namespaceQuery := "oci_postgresql"
	query := "Connections[1m].max()"

	compartmentId := config.CompartmentId

	err := getPostgresqlConnectionsByCompartment(
		ctx,
		compartmentId,
		query,
		namespaceQuery,
		postgresqlConnections,
	)
	if err != nil {
		return err
	}

	return nil
}

func getPostgresqlConnectionsByCompartment(
	ctx context.Context,
	compartmentId string,
	query string,
	namespaceQuery string,
	postgresqlConnections *prometheus.GaugeVec,
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
		postgresqlConnections.With(prometheus.Labels{
			"resource_name":    resourceName,
			"resource_id":      resourceId,
			"compartment_id":   compartmentId,
			"db_instance_id":   dbInstanceId,
			"db_instance_role": dbInstanceRole,
		}).Set(value)
	}

	return nil
}
