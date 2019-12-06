package telphin

import "fmt"

// Endpoint: GET /extension/{extension_id}/ani/
func (c *Client) GetExtensionAni(extensionID uint16) (*Ani, error) {
	ani := &Ani{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s/extension/%d/ani", c.Host, extensionID), nil)
	if err != nil {
		return ani, err
	}

	if err = c.SendWithAuth(req, ani); err != nil {
		return ani, err
	}
	return ani, nil
}

// Endpoint: PUT /extension/{extension_id}/ani/
func (c *Client) SetExtensionAni(extensionID uint32, aniNumber string) (*Ani, error) {
	type aniRequest struct {
		AniNumber string `json:"ani_number"`
	}

	ani := &Ani{}

	req, err := c.NewRequest("PUT", fmt.Sprintf("%s/api/ver1.0/extension/%d/ani/", c.Host, extensionID), aniRequest{AniNumber: aniNumber})
	if err != nil {
		return ani, err
	}

	if err = c.SendWithAuth(req, ani); err != nil {
		return ani, err
	}
	return ani, nil
}
