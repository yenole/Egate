//+build !dws,client

package network

import (
	"egate/gate"
	"fmt"
	"github.com/gorilla/websocket"
)

type wsClientAgent struct {
	*WsConn
	clt *gate.Client
}

func (t *wsClientAgent) AgentWriteMsg(m interface{}) {
	t.clt.Answer(t, m)
}

func (t *wsClientAgent) Close() error {
	return t.Conn.Close()
}

func NewWsClient(addr string, client *gate.Client) *wsClientAgent {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%v/", addr), nil)
	if err != nil {
		panic(err)
	}
	return &wsClientAgent{WsConn: NewWsConn(conn), clt: client}
}
