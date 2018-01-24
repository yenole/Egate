package network

import (
	"egate/elog"
	"egate/gate"
	"net"
	"sync"
)

type TcpNetworkAccpet struct {
	ln net.Listener // listen
}

func NewTcpNetworkAccpet() *TcpNetworkAccpet {
	return &TcpNetworkAccpet{}
}

// init tcp
func (t *TcpNetworkAccpet) Init(opt gate.Config) bool {
	cfg := Config(opt)
	if len(cfg.address()) > 0 {
		ln, err := net.Listen("tcp", cfg.address())
		if err != nil {
			panic(err)
		}
		t.ln = ln
		elog.Info("tcp accpet listen[%v]", cfg.address())
		return true
	}
	return false
}

// Work
func (t *TcpNetworkAccpet) Work(wg *sync.WaitGroup, egt gate.RecvAnswer) {
	if wg != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				cn, err := t.ln.Accept()
				if err != nil {
					elog.Info("%v", err)
					break
				}
				egt.Recv(&tcpAgent{Conn: cn, tcp: t, egt: egt})
			}
		}()
	}
}

type tcpAgent struct {
	net.Conn
	tcp *TcpNetworkAccpet
	egt gate.RecvAnswer
}

func (t *tcpAgent) AgentWriteMsg(m interface{}) {
	t.egt.Answer(t, m)
}

func (t *tcpAgent) Close() error {
	return t.Conn.Close()
}
