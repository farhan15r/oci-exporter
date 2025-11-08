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

func GetInstanceCPUUtilization(ctx context.Context) (string, error) {
	metricName := "oci_exporter_instance_cpu_utilization"
	result := "# HELP oci_exporter_instance_cpu_utilization CPU Utilization of OCI Instance, unit Percent\n"
	result += "# TYPE oci_exporter_instance_cpu_utilization gauge\n"

	compartmentIds := config.CompartmentIds

	for _, compartmentId := range compartmentIds {
		metricsData, err := getCPUUtilByCompartment(ctx, compartmentId, metricName)
		if err != nil {
			return "", err
		}
		result += metricsData
	}

	return result, nil
}

func getCPUUtilByCompartment(ctx context.Context, compartmentId string, metricName string) (string, error) {
	minutes := 1

	end := time.Now().UTC()
	start := end.Add(-time.Duration(minutes) * time.Minute)

	query := "CPUUtilization[1m].mean()"
	namespace := "oci_computeagent"
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
		value := *lastPoint.Value

		// extract dimension values
		instanceName := metric.Dimensions["resourceDisplayName"]
		compartmentID := *metric.CompartmentId
		faultDomain := metric.Dimensions["faultDomain"]
		resourceID := metric.Dimensions["resourceId"]

		// format result line
		result += fmt.Sprintf("%s{instance_name=\"%s\", compartment_id=\"%s\", fault_domain=\"%s\", resource_id=\"%s\"} %f\n",
			metricName,
			instanceName,
			compartmentID,
			faultDomain,
			resourceID,
			value,
		)
	}

	return result, nil
}
