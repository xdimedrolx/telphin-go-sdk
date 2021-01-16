package telphin

import "fmt"

// Endpoint: GET /extension/{extension_id}/event/
func (c *Client) GetEvents(extensionId uint32) (*[]Event, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/extension/%d/event/", c.Host, extensionId), nil)
	if err != nil {
		return nil, err
	}

	events := &[]Event{}
	if err = c.SendWithAuth(req, events); err != nil {
		return nil, err
	}
	return events, nil
}

// Endpoint: POST /extension/{extension_id}/event/
func (c *Client) CreateEvent(extensionId uint32, event CreateEventRequest) (*Event, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s/api/ver1.0/extension/%d/event/", c.Host, extensionId), event)
	if err != nil {
		return nil, err
	}

	resp := &Event{}
	if err = c.SendWithAuth(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Endpoint: GET /extension/{extension_id}/event/{id}
func (c *Client) DeleteEvent(extensionId uint32, eventId int) error {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("%s/api/ver1.0/extension/%d/event/%d", c.Host, extensionId, eventId), nil)
	if err != nil {
		return err
	}

	return c.SendWithAuth(req, nil)
}
