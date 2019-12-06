package telphin

import "fmt"

// Endpoint: PUT /extension/{extension_id}/phone
func (c *Client) UpdatePhone(extensionID uint32, phoneProps PhoneProperties) error {
	req, err := c.NewRequest("PUT", fmt.Sprintf("%s/api/ver1.0/extension/%d/phone/", c.Host, extensionID), phoneProps)
	if err != nil {
		return err
	}

	newPhoneProps := &PhoneProperties{}
	if err = c.SendWithAuth(req, newPhoneProps); err != nil {
		return err
	}
	return nil
}
