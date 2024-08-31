package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/rs/zerolog/log"
)

func (c *Client) GetGateways() ([]Gateway, error) {
	devices, err := c.GetDevices()
	if err != nil {
		return nil, err
	}

	gateways := []Gateway{}
	for _, d := range devices {
		if d.Type == "gateway" {
			gateway, err := c.GetGateway(d.Mac)
			if err != nil {
				return nil, fmt.Errorf("failed to get gateway: %s", err)
			}
			gateways = append(gateways, *gateway)
		}
	}

	return gateways, nil
}

func (c *Client) GetGateway(switchMac string) (*Gateway, error) {
	url := fmt.Sprintf("%s/%s/api/v2/sites/%s/gateways/%s", c.Config.Host, c.omadaCID, c.SiteId, switchMac)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeLoggedInRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debug().Bytes("data", body).Msg("Received data from ports endpoint")

	portdata := gatewayResponse{}
	err = json.Unmarshal(body, &portdata)

	return portdata.Result, err
}

type gatewayResponse struct {
	Result *Gateway `json:"result"`
}

type Gateway struct {
	Mac       string             `json:"mac"`
	Name      string             `json:"name"`
	Model     string             `json:"model"`
	PortStats []GatewayPortStats `json:"portStats"`
}

type GatewayPortStats struct {
	Port            int64   `json:"port"`
	Type            int64   `json:"type"`
	Mode            int64   `json:"mode"`
	Name            string  `json:"name"`
	InternetState   float64 `json:"internetState"`
	Description     string  `json:"portDesc"`
	OnlineDetection int64   `json:"onlineDetection" default:"0"`
}
