package telphin

import (
	"fmt"

	mem "github.com/spf13/afero/mem"
)

// Endpoint: GET /client/{client_id}/record/{record_uuid}/storage_url/
func (c *Client) GetRecordStorageUrl(clientID string, recordUUID string) (*RecordStorageUrl, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/record/%s/storage_url/", c.Host, clientID, recordUUID), nil)
	if err != nil {
		return nil, err
	}

	record := &RecordStorageUrl{}
	if err = c.SendWithAuth(req, record); err != nil {
		return nil, err
	}
	return record, nil
}

// Endpoint: GET /client/{client_id}/record/{record_uuid}
func (c *Client) GetRecord(clientID string, recordUUID string) (*mem.File, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/record/%s", c.Host, clientID, recordUUID), nil)
	if err != nil {
		return nil, err
	}

	file := mem.NewFileHandle(mem.CreateFile("record.mp3"))
	if err = c.SendWithAuth(req, file); err != nil {
		return nil, err
	}
	return file, nil
}
