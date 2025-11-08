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

var fastconnectBgpSession = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "oci_exporter",
		Name:      "fastconnect_ipv4_bgp_session_state",
		Help:      "BGP State of OCI IPv4 FastConnect, 1=up, 0=down.",
	},
	[]string{"resource_name", "compartment_id", "resource_id"},
)

func GetFastconnectBGPSessionState(ctx context.Context) (*prometheus.GaugeVec, error) {
	fastconnectBgpSession.Reset()

	namespaceQuery := "oci_fastconnect"
	query := "Ipv4BgpSessionState[1m].mean()"

	compartmentId := config.CompartmentId

	err := getFastconnectBGPSessionStateByCompartment(
		ctx,
		fastconnectBgpSession,
		compartmentId,
		query,
		namespaceQuery,
	)
	if err != nil {
		return nil, err
	}

	return fastconnectBgpSession, nil
}

func getFastconnectBGPSessionStateByCompartment(
	ctx context.Context,
	fastconnectBgpSession *prometheus.GaugeVec,
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
		resourceName := metric.Dimensions["resourceName"]
		compartmentId := *metric.CompartmentId
		resourceId := metric.Dimensions["resourceId"]

		// set gauge value
		fastconnectBgpSession.With(prometheus.Labels{
			"resource_name":  resourceName,
			"compartment_id": compartmentId,
			"resource_id":    resourceId,
		}).Set(float64(value))
	}

	return nil
}
