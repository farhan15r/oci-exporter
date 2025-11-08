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

var vpnBytesSent = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "oci_exporter",
		Name:      "vpn_bytes_sent",
		Help:      "Total Bytes Sent on OCI VPN.",
	},
	[]string{"resource_name", "compartment_id", "parent_resource_id"},
)

func GetVpnBytesSentState(ctx context.Context) (*prometheus.GaugeVec, error) {
	vpnBytesSent.Reset()

	namespaceQuery := "oci_vpn"
	query := "BytesSent[1m].sum()"

	compartmentId := config.CompartmentId

	err := getVpnBytesSentStateByCompartment(
		ctx,
		compartmentId,
		query,
		namespaceQuery,
	)
	if err != nil {
		return nil, err
	}

	return vpnBytesSent, nil
}

func getVpnBytesSentStateByCompartment(
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
		compartmentId := *metric.CompartmentId
		parentResourceId := metric.Dimensions["parentResourceId"]

		// set gauge value
		vpnBytesSent.With(prometheus.Labels{
			"resource_name":      resourceName,
			"compartment_id":     compartmentId,
			"parent_resource_id": parentResourceId,
		}).Set(value)
	}

	return nil
}
