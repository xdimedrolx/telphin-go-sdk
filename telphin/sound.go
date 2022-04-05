package telphin

import (
	"fmt"
	"github.com/google/go-querystring/query"
)

// Endpoint: GET /client/{client_id}/sound/
func (c *Client) GetSounds(clientID string, request SoundsRequest) (*[]Sound, error) {
	sounds := &[]Sound{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/sound/", c.Host, clientID), nil)
	if err != nil {
		return sounds, err
	}
	v, _ := query.Values(request)
	req.URL.RawQuery = v.Encode()

	if err = c.SendWithAuth(req, sounds); err != nil {
		return nil, err
	}
	return sounds, nil
}
