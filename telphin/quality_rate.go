package telphin

import (
	"fmt"
)

// GET /api/ver1.0/client/{client_id}/quality_rate/
func (c *Client) GetQualityRate(clientID string, request QualityRateRequest) (*[]QualityRate, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/quality_rate/", c.Host, clientID), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	// FIXME
	if request.ExtensionID != nil {
		q.Add("extension_id", fmt.Sprint(*request.ExtensionID))
	}
	if request.StartDatetime != nil {
		q.Add("start_datetime", request.StartDatetime.Format("2006-01-02 15:04:05"))
	}
	if request.EndDatetime != nil {
		q.Add("end_datetime", request.EndDatetime.Format("2006-01-02 15:04:05"))
	}
	req.URL.RawQuery = q.Encode()

	rates := &[]QualityRate{}
	if err = c.SendWithAuth(req, rates); err != nil {
		return nil, err
	}
	return rates, nil
}
