package telphin

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

// GET /api/ver1.0/client/{client_id}/quality_rate/
func (c *Client) GetQualityRate(clientID string, request QualityRateRequest) (*[]QualityRate, error) {
	rates := &[]QualityRate{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/quality_rate/", c.Host, clientID), nil)
	if err != nil {
		return rates, err
	}
	v, _ := query.Values(request)
	req.URL.RawQuery = v.Encode()

	if err = c.SendWithAuth(req, rates); err != nil {
		return rates, err
	}
	return rates, nil
}
