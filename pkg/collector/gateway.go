package collector

import (
	"github.com/charlie-haley/omada_exporter/pkg/api"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/rs/zerolog/log"
)

type gatewayCollector struct {
	gatewayWanPortInternetOnline *prometheus.Desc
	client                       *api.Client
}

func (c *gatewayCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (c *gatewayCollector) Collect(ch chan<- prometheus.Metric) {
	client := c.client
	config := c.client.Config

	site := config.Site
	gateways, err := client.GetGateways()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get gateways")
		return
	}

	for _, item := range gateways {
		for _, port := range item.PortStats {
			if port.Type <= 1 && port.Mode == 0 {
				labels := []string{item.Name, item.Mac, item.Model, site, client.SiteId, port.Name}
				ch <- prometheus.MustNewConstMetric(c.gatewayWanPortInternetOnline,
					prometheus.GaugeValue, float64(port.OnlineDetection), labels...)
			}
		}
	}
}

func NewGatewayCollector(c *api.Client) *gatewayCollector {
	labels := []string{"name", "mac", "model", "site", "site_id", "port"}

	return &gatewayCollector{
		gatewayWanPortInternetOnline: prometheus.NewDesc("omada_gateway_wan_port_internet_online",
			"Internet active state of the WAN port", labels, nil),
		client: c,
	}
}
