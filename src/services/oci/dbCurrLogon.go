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

var dbCurrLogon = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "oci_exporter",
		Name:      "database_current_logons",
		Help:      "The number of successful logons",
	},
	[]string{
		"reource_id",
		"compartment_id",
		"resource_name",
		"instance_name",
	},
)

func GetDbCurrLogon(ctx context.Context) (*prometheus.GaugeVec, error) {
	dbCurrLogon.Reset()

	namespaceQuery := "oci_database"
	query := "CurrentLogons[1m].mean()"

	compartmentId := config.CompartmentId

	err := getDbCurrLogonByCompartment(
		ctx,
		compartmentId,
		query,
		namespaceQuery,
	)
	if err != nil {
		return nil, err
	}

	return dbCurrLogon, nil
}

func getDbCurrLogonByCompartment(
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
		resourceName := metric.Dimensions["resourceName"]
		instanceName := metric.Dimensions["instanceName"]

		compartmentId := *metric.CompartmentId

		// set gauge value
		dbCurrLogon.With(prometheus.Labels{
			"reource_id":     resourceId,
			"compartment_id": compartmentId,
			"resource_name":  resourceName,
			"instance_name":  instanceName,
		}).Set(value)
	}

	return nil
}
