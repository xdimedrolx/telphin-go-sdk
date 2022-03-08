package telphin

import (
	"fmt"
	"github.com/google/go-querystring/query"

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

// Endpoint: GET /client/{client_id}/record/
func (c *Client) GetRecords(clientID string, recordsRequest RecordsRequest) (*[]RecordInfo, error) {
	records := &[]RecordInfo{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/record/", c.Host, clientID), nil)
	if err != nil {
		return records, err
	}
	v, _ := query.Values(recordsRequest)
	req.URL.RawQuery = v.Encode()

	if err = c.SendWithAuth(req, records); err != nil {
		return nil, err
	}
	return records, nil
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

// Endpoint: DELETE /client/{client_id}/record/{record_uuid}
func (c *Client) DeleteRecord(clientID string, recordUUID string) error {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("%s/api/ver1.0/client/%s/record/%s", c.Host, clientID, recordUUID), nil)
	if err != nil {
		return err
	}

	return c.SendWithAuth(req, nil)
}
