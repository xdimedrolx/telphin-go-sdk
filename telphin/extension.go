package telphin

import (
	"fmt"
	"net/url"
)

// Endpoint: GET /client/{client_id}/extension/
func (c *Client) GetExtensions(clientID string, typeExtensions *string, page *int) (*[]Extension, error) {
	extensions := &[]Extension{}

	u, _ := url.Parse(fmt.Sprintf("%s/api/ver1.0/client/%s/extension/", c.Host, clientID))

	q := url.Values{}
	if typeExtensions != nil {
		q.Add("type", *typeExtensions)
	}
	if page != nil {
		q.Add("page", string(*page))
	}

	u.RawQuery = q.Encode()

	req, err := c.NewRequest("GET", u.String(), nil)
	if err != nil {
		return extensions, err
	}

	if err = c.SendWithAuth(req, extensions); err != nil {
		return extensions, err
	}

	return extensions, nil
}

// Endpoint: GET /client/{client_id}/extension/{extension_id}
func (c *Client) GetExtension(clientID string, extensionID uint16) (*Extension, error) {
	extension := &Extension{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/extension/%d", c.Host, clientID, extensionID), nil)
	if err != nil {
		return extension, err
	}

	if err = c.SendWithAuth(req, extension); err != nil {
		return extension, err
	}
	return extension, nil
}

// Endpoint: POST /client/{client_id}/extension
func (c *Client) CreateExtension(clientID string, extensionCreateRequest ExtensionCreateRequest) (*Extension, error) {
	extension := &Extension{}
	req, err := c.NewRequest("POST", fmt.Sprintf("%s/api/ver1.0/client/%s/extension/", c.Host, clientID), extensionCreateRequest)
	if err != nil {
		return extension, err
	}

	if err = c.SendWithAuth(req, extension); err != nil {
		return extension, err
	}
	return extension, nil
}

// Endpoint: DELETE /client/{client_id}/extension/{extension_id}
func (c *Client) DeleteExtension(clientID string, extensionID uint32) error {
	req, err := c.NewRequest("DELETE", fmt.Sprintf("%s/api/ver1.0/client/%s/extension/%d", c.Host, clientID, extensionID), nil)
	if err != nil {
		return err
	}

	err = c.SendWithAuth(req, nil)
	return err
}

// Endpoint: POST /extension/{extension_id}/callback/
func (c *Client) CreateCallback(extensionID uint32, callback CallbackRequest) (*Callback, error) {
	req, err := c.NewRequest("POST", fmt.Sprintf("%s/api/ver1.0/extension/%d/callback", c.Host, extensionID), callback)
	if err != nil {
		return nil, err
	}

	resp := &Callback{}
	err = c.SendWithAuth(req, resp)
	return resp, err
}
