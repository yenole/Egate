package gate

import (
	"io"
	"net"
	"sync"
)

type Config map[int]interface{}

type MsgInfo struct {
	Id  uint32
	Msg interface{}
}

type Accpeter interface {
	Init(opt Config) bool
	Work(wg *sync.WaitGroup, egt RecvAnswer)
}

type RecvAnswer interface {
	Recv(agent Agent)
	Answer(agent Agent, m interface{})
}

type Agent interface {
	io.ReadWriteCloser
	AgentWriteMsg(m interface{})
	RemoteAddr() net.Addr
}

type IDer interface {
	ID() uint64
	SetId(v uint64)
}
