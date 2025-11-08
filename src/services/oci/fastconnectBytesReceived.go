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

var fastconnectBytesReceived = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "oci_exporter",
		Name:      "fastconnect_bytes_received",
		Help:      "Total Bytes Received on OCI FastConnect.",
	},
	[]string{"resource_name", "compartment_id", "resource_id"},
)

func GetFastconnectBytesReceived(ctx context.Context) (*prometheus.GaugeVec, error) {
	fastconnectBytesReceived.Reset()

	namespaceQuery := "oci_fastconnect"
	query := "BytesReceived[1m].sum()"

	compartmentId := config.CompartmentId

	err := getFastconnectBytesReceivedByCompartment(
		ctx,
		compartmentId,
		query,
		namespaceQuery,
	)
	if err != nil {
		return nil, err
	}

	return fastconnectBytesReceived, nil
}

func getFastconnectBytesReceivedByCompartment(
	ctx context.Context,
	compartmentId string,
	query string,
	namespaceQuery string,
) error {
	minutes := 5

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
		compartmentID := *metric.CompartmentId
		resourceID := metric.Dimensions["resourceId"]

		fastconnectBytesReceived.With(prometheus.Labels{
			"resource_name":  resourceName,
			"compartment_id": compartmentID,
			"resource_id":    resourceID,
		}).Set(value)
	}

	return nil
}
