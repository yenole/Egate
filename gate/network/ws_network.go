//+build !dws

package network

import (
	"egate/elog"
	"egate/gate"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type WsNetworkAccpet struct {
	egt      gate.RecvAnswer
	server   *http.Server
	upgrader websocket.Upgrader
}

func NewWsNetworkAccpet() *WsNetworkAccpet {
	return &WsNetworkAccpet{}
}

// init ws
func (t *WsNetworkAccpet) Init(opt gate.Config) bool {
	cfg := Config(opt)
	if len(cfg.address()) > 0 {
		t.server = &http.Server{Addr: cfg.address(), Handler: t}
		elog.Info("ws accpet listen[%v]", cfg.address())
		return true
	}
	return false
}

// Work
func (t *WsNetworkAccpet) Work(wg *sync.WaitGroup, egt gate.RecvAnswer) {
	t.egt = egt
	if t.server != nil && wg != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.server.ListenAndServe()
		}()
	}
}

func (t *WsNetworkAccpet) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	req.Header.Del("Origin")
	conn, err := t.upgrader.Upgrade(resp, req, nil)
	if err != nil {
		elog.Info("%v", err)
		return
	}
	t.egt.Recv(&wsAgent{WsConn: NewWsConn(conn), ws: t, egt: t.egt})
}

type wsAgent struct {
	*WsConn
	ws  *WsNetworkAccpet
	egt gate.RecvAnswer
}

func (t *wsAgent) AgentWriteMsg(m interface{}) {
	t.egt.Answer(t, m)
}

func (t *wsAgent) Close() error {
	return t.Conn.Close()
}
