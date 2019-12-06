package telphin

import "fmt"

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
