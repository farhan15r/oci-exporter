package oci

import (
	"context"
	"fmt"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/monitoring"

	"oci-exporter/src/config"
	"oci-exporter/src/utils"
)

func GetFastconnectBytesReceived(ctx context.Context) (string, error) {
	// map[string]string ["displayName": "Bytes Received", "unit": "bytes", ]
	metricName := "oci_exporter_fastconnect_bytes_received"
	result := "# HELP oci_exporter_fastconnect_bytes_received unit: bytes. FastConnect Bytes Received.\n"
	result += "# TYPE oci_exporter_fastconnect_bytes_received gauge.\n"

	compartmentIds := config.CompartmentIds

	for _, compartmentId := range compartmentIds {
		metricsData, err := getFastconnectBytesReceivedByCompartment(ctx, compartmentId, metricName)
		if err != nil {
			return "", err
		}
		result += metricsData
	}

	return result, nil
}

func getFastconnectBytesReceivedByCompartment(ctx context.Context, compartmentId string, metricName string) (string, error) {
	minutes := 5

	end := time.Now().UTC()
	start := end.Add(-time.Duration(minutes) * time.Minute)

	query := "BytesReceived[1m].sum()"
	namespace := "oci_fastconnect"
	sdkStart := common.SDKTime{Time: start}
	sdkEnd := common.SDKTime{Time: end}

	req := monitoring.SummarizeMetricsDataRequest{
		CompartmentId: &compartmentId,
		SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
			Query:     &query,
			StartTime: &sdkStart,
			EndTime:   &sdkEnd,
			Namespace: &namespace,
		},
	}

	client, err := config.NewOciClient()
	if err != nil {
		utils.Logger.Error("failed to create OCI client", "error", err.Error())
		return "", err
	}

	resp, err := client.SummarizeMetricsData(ctx, req)
	if err != nil {
		utils.Logger.Error("SummarizeMetricsData failed", "error", err.Error())
		return "", err
	}

	var result string

	for _, metric := range resp.Items {
		// get last data point
		if len(metric.AggregatedDatapoints) == 0 {
			continue
		}
		lastPoint := metric.AggregatedDatapoints[len(metric.AggregatedDatapoints)-1]
		value := int(*lastPoint.Value)

		// extract dimension values
		resourceName := metric.Dimensions["resourceName"]
		compartmentID := *metric.CompartmentId
		resourceID := metric.Dimensions["resourceId"]

		// format result line
		result += fmt.Sprintf("%s{resourceName=\"%s\", compartment_id=\"%s\", resource_id=\"%s\"} %d\n",
			metricName,
			resourceName,
			compartmentID,
			resourceID,
			value,
		)
	}

	return result, nil
}
