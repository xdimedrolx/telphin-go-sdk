package telphin

import (
	"fmt"
	"net/url"
)

// Endpoint: GET /client/{client_id}/extension/
func (c *Client) GetExtensions(clientID string, typeExtensions *string, page *int) (*[]Extension, error) {
	u, _ := url.Parse(fmt.Sprintf("%s/api/ver1.0/client/%s/extension/", c.Host, clientID))

	q := url.Values{}
	if typeExtensions != nil {
		q.Add("type", *typeExtensions)
	}
	if page != nil {
		q.Add("page", fmt.Sprint(*page))
	}

	u.RawQuery = q.Encode()

	req, err := c.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	extensions := &[]Extension{}
	if err = c.SendWithAuth(req, extensions); err != nil {
		return nil, err
	}

	return extensions, nil
}

// Endpoint: GET /client/{client_id}/extension/{extension_id}
func (c *Client) GetExtension(clientID string, extensionID uint16) (*Extension, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/extension/%d", c.Host, clientID, extensionID), nil)
	if err != nil {
		return nil, err
	}

	extension := &Extension{}
	if err = c.SendWithAuth(req, extension); err != nil {
		return nil, err
	}
	return extension, nil
}

// Endpoint: POST /client/{client_id}/extension
func (c *Client) CreateExtension(clientID string, extensionCreateRequest ExtensionCreateRequest) (*Extension, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s/api/ver1.0/client/%s/extension/", c.Host, clientID), extensionCreateRequest)
	if err != nil {
		return nil, err
	}

	extension := &Extension{}
	if err = c.SendWithAuth(req, extension); err != nil {
		return nil, err
	}
	return extension, nil
}

// Endpoint: DELETE /client/{client_id}/extension/{extension_id}
func (c *Client) DeleteExtension(clientID string, extensionID uint32) error {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("%s/api/ver1.0/client/%s/extension/%d", c.Host, clientID, extensionID), nil)
	if err != nil {
		return err
	}

	return c.SendWithAuth(req, nil)
}

// Endpoint: POST /api/ver1.0/extension/{extension_id}/callback/
func (c *Client) CreateCallback(extensionID uint32, callback CallbackRequest) (*Callback, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s/api/ver1.0/extension/%d/callback/", c.Host, extensionID), callback)
	if err != nil {
		return nil, err
	}

	resp := &Callback{}
	if err = c.SendWithAuth(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
