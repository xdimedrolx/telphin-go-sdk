package telphin

import (
	"fmt"
	"net/http"
)

// Endpoint: GET /client/{client_id}/record/{record_uuid}/storage_url/
func (c *Client) GetRecordStorageUrl(clientID string, recordUUID string) (*RecordStorageUrl, error) {
	record := &RecordStorageUrl{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/record/%s/storage_url/", c.Host, clientID, recordUUID), nil)
	if err != nil {
		return record, err
	}
	err = c.SendWithAuth(req, record)
	return record, err
}

// Endpoint: GET /client/{client_id}/record/{record_uuid}
func (c *Client) GetRecord(clientID string, recordUUID string) (*http.Response, error) {
	record := &http.Response{}
	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/record/%s", c.Host, clientID, recordUUID), nil)
	if err != nil {
		return nil, err
	}
	err = c.SendWithAuth(req, record)
	return record, err
}
