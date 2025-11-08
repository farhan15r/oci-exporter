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

var vpnIpSecState = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "oci_exporter",
		Name:      "vpn_ipsec_tunnel_state",
		Help:      "IPSec Tunnel State of OCI VPN, 1=up, 0=down.",
	},
	[]string{"resource_name", "compartment_id", "parent_resource_id"},
)

func GetVpnIpSecStateState(ctx context.Context) (*prometheus.GaugeVec, error) {
	vpnIpSecState.Reset()

	namespaceQuery := "oci_vpn"
	query := "TunnelState[1m].mean()"

	// create local GaugeVec and return it; caller (handler) should register if desired

	compartmentId := config.CompartmentId

	err := getVpnIpSecStateStateByCompartment(
		ctx,
		compartmentId,
		query,
		namespaceQuery,
	)
	if err != nil {
		return nil, err
	}

	return vpnIpSecState, nil
}

func getVpnIpSecStateStateByCompartment(
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
		resourceName := metric.Dimensions["resourceName"]
		compartmentId := *metric.CompartmentId
		parentResourceId := metric.Dimensions["parentResourceId"]

		// set gauge value
		vpnIpSecState.With(prometheus.Labels{
			"resource_name":      resourceName,
			"compartment_id":     compartmentId,
			"parent_resource_id": parentResourceId,
		}).Set(float64(value))
	}

	return nil
}
