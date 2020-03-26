package telphin

import "fmt"

func (c *Client) GetAllDID(clientID string) (*[]Did, error) {
	dids := &[]Did{}

	req, err := c.NewRequest("GET", fmt.Sprintf("%s/api/ver1.0/client/%s/did", c.Host, clientID), nil)
	if err != nil {
		return dids, err
	}

	if err = c.SendWithAuth(req, dids); err != nil {
		return dids, err
	}
	return dids, nil
}
