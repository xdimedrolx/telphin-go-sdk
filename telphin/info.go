package telphin

func (c *Client) WsHost() string {
	return WsHost
}

// WsServers is WebSocket URIs to connect to.
func (c *Client) WsServers() []string {
	return []string{"wss://sipproxy.telphin.ru", "wss://pbx.telphin.ru"}
}
