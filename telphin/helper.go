package telphin

func (c *Client) GetWsSipUri(user string) string {
	return user + "@" + WsHost
}
