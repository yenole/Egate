//+build client

package gate

type Client struct {
	// Agents
	agent  Agent
	mfsin  []MiddlewareFunc
	mfsout []MiddlewareFunc
}

func NewClient() *Client {
	return &Client{
		mfsin:  make([]MiddlewareFunc, 0),
		mfsout: make([]MiddlewareFunc, 0),
	}
}

func (c *Client) WriteAgentMsg(m interface{}) {
	c.agent.AgentWriteMsg(m)
}

func (c *Client) Agent(agent Agent) *Client {
	c.agent = agent
	return c
}

func (c *Client) In(f MiddlewareFunc) *Client {
	if f != nil {
		c.mfsin = append(c.mfsin, f)
	}
	return c
}

func (c *Client) Out(f MiddlewareFunc) *Client {
	if f != nil {
		c.mfsout = append(c.mfsout, f)
	}
	return c
}

func (c *Client) Answer(agent Agent, m interface{}) {
	p := Middleware{mfs: c.mfsout, Agent: agent}
	p.Msg.Msg = m
	p.Next()
}

func (c *Client) Recv() {
	m := Middleware{mfs: c.mfsin, Agent: c.agent}
	for {
		m.next = 0
		if m.Next().IsAbort() {
			break
		}
	}
}
