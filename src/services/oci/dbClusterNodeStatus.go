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

var dbClusterNodeStatus = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "oci_exporter",
		Name:      "database_cluster_node_status",
		Help:      "Node Status of OCI Database Cluster. 1 for Up, 0 for Down.",
	},
	[]string{
		"compartment_id",
		"db_node_name",
		"reource_id",
		"resource_name",
	},
)

func GetDbClusterNodeStatus(ctx context.Context) (*prometheus.GaugeVec, error) {
	dbClusterNodeStatus.Reset()

	namespaceQuery := "oci_database_cluster"
	query := "NodeStatus[1m].mean()"

	compartmentId := config.CompartmentId

	err := getDbClusterNodeStatusByCompartment(
		ctx,
		compartmentId,
		query,
		namespaceQuery,
	)
	if err != nil {
		return nil, err
	}

	return dbClusterNodeStatus, nil
}

func getDbClusterNodeStatusByCompartment(
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
		value := int(*lastPoint.Value)

		// extract dimension values
		compartmentId := *metric.CompartmentId
		dbNodeName := metric.Dimensions["resourceName_dbnode"]
		resourceId := metric.Dimensions["resourceId"]
		resourceName := metric.Dimensions["resourceName"]

		// set gauge value
		dbClusterNodeStatus.With(prometheus.Labels{
			"compartment_id": compartmentId,
			"db_node_name":   dbNodeName,
			"reource_id":     resourceId,
			"resource_name":  resourceName,
		}).Set(float64(value))
	}

	return nil
}
