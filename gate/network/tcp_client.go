//+build client

package network

import (
	"egate/gate"
	"net"
)

type tcpClientAgent struct {
	net.Conn
	clt *gate.Client
}

func (t *tcpClientAgent) AgentWriteMsg(m interface{}) {
	t.clt.Answer(t, m)
}

func (t *tcpClientAgent) Close() error {
	return t.Conn.Close()
}

func NewTCPClient(addr string, client *gate.Client) *tcpClientAgent {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	return &tcpClientAgent{Conn: conn, clt: client}
}
