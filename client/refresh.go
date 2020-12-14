package client

func (c *client) Refresh() {
	c.sendStringMessage("refresh", "")
}
