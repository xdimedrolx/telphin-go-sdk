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
	err = c.SendWithAuth(req, record)
	return record, err
}

// Endpoint: GET /client/{client_id}/record/{record_uuid}
func (c *Client) GetRecord(clientID string, recordUUID string) (*mem.File, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/record/%s", c.Host, clientID, recordUUID), nil)
	if err != nil {
		return nil, err
	}

	file := mem.NewFileHandle(mem.CreateFile("record.mp3"))
	err = c.SendWithAuth(req, file)
	return file, err
}
