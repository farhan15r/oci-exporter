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

func newPostgresqlUsedStorage() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "oci_exporter",
			Name:      "postgresql_used_storage",
			Help:      "Used Storage of OCI PostgreSQL Database.",
		},
		[]string{
			"resource_name",
			"resource_id",
			"compartment_id",
			"db_instance_id",
			"resource_type",
		},
	)
}

func GetPostgresqlUsedStorage(ctx context.Context, postgresqlUsedStorage *prometheus.GaugeVec) error {
	postgresqlUsedStorage.Reset()

	namespaceQuery := "oci_postgresql"
	query := "UsedStorage[1m].max()"

	compartmentId := config.CompartmentId

	err := getPostgresqlUsedStorageByCompartment(
		ctx,
		compartmentId,
		query,
		namespaceQuery,
		postgresqlUsedStorage,
	)
	if err != nil {
		return err
	}

	return nil
}

func getPostgresqlUsedStorageByCompartment(
	ctx context.Context,
	compartmentId string,
	query string,
	namespaceQuery string,
	postgresqlUsedStorage *prometheus.GaugeVec,
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
		resourceType := metric.Dimensions["resourceType"]

		// set gauge value
		postgresqlUsedStorage.With(prometheus.Labels{
			"resource_name":  resourceName,
			"resource_id":    resourceId,
			"compartment_id": compartmentId,
			"db_instance_id": dbInstanceId,
			"resource_type":  resourceType,
		}).Set(value)
	}

	return nil
}
