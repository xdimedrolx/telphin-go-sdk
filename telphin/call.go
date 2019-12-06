package telphin

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

// Endpoint: GET /client/{client_id}/call_history/
func (c *Client) GetCallHistories(clientID string, callHistoryRequest CallHistoryRequest) (*CallHistories, error) {
	histories := &CallHistories{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/call_history/", c.Host, clientID), nil)
	if err != nil {
		return histories, err
	}
	v, _ := query.Values(callHistoryRequest)
	req.URL.RawQuery = v.Encode()

	err = c.SendWithAuth(req, histories)
	return histories, err
}

func (c *Client) GetCallHistory(clientID string, callId string) (*CallHistory, error) {
	call := &CallHistory{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/call_history/%s", c.Host, clientID, callId), nil)
	if err != nil {
		return call, err
	}

	err = c.SendWithAuth(req, call)
	return call, err
}
