package gate

import (
	"egate/assist"
	"egate/elog"
	"runtime"
	"sync"
)

func Version() string {
	return "0.1.0"
}

type MiddlewareFunc func(m *Middleware)

type Egate struct {
	accepts *assist.EList
	plugins *assist.EMap

	// Agents
	agents *assist.EMap
	mfsin  []MiddlewareFunc
	mfsout []MiddlewareFunc
}

func NewEgate() *Egate {
	return &Egate{
		accepts: assist.NewElist(),
		plugins: assist.NewEMap(),
		agents:  assist.NewEMap(),
		mfsin:   make([]MiddlewareFunc, 0),
		mfsout:  make([]MiddlewareFunc, 0),
	}
}

func (e *Egate) In(f MiddlewareFunc) {
	if f != nil {
		e.mfsin = append(e.mfsin, f)
	}
}

func (e *Egate) Out(f MiddlewareFunc) *Egate {
	if f != nil {
		e.mfsout = append(e.mfsout, f)
	}
	return e
}

func (e *Egate) Accpet(a Accpeter) *Egate {
	if a != nil {
		e.accepts.Push(a)
	}
	return e
}

func (e *Egate) Work(wg *sync.WaitGroup) bool {
	if wg != nil {
		e.accepts.ForEach(func(v interface{}, i int) bool {
			v.(Accpeter).Work(wg, e)
			return true
		})
		return true
	}
	return false
}

func (e *Egate) Answer(agent Agent, m interface{}) {
	p := Middleware{mfs: e.mfsout, Agent: agent}
	p.Msg.Msg = m
	p.Next()
}

func (e *Egate) Recv(agent Agent) {
	e.agents.Push(agent, true)
	go func() {
		defer agent.Close()
		defer e.agents.Pop(agent)
		defer func() {
			if err := recover(); err != nil {
				elog.Fatal("%v", err)
				elog.Fatal("%v", string(Stack()))
			}
		}()

		m := Middleware{mfs: e.mfsin, Agent: agent}
		for {
			m.next = 0
			if m.Next().IsAbort() {
				break
			}
		}
	}()
}

type Middleware struct {
	Agent Agent
	Msg   MsgInfo

	abort bool
	next  int
	mfs   []MiddlewareFunc
	extra map[uint32]interface{}
}

func (m *Middleware) Extra() map[uint32]interface{} {
	if m.extra == nil {
		m.extra = make(map[uint32]interface{})
	}
	return m.extra
}

func (m *Middleware) Abort() {
	m.abort = true
}

func (m *Middleware) IsAbort() bool {
	return m.abort
}

func Stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}

func (m *Middleware) Next() *Middleware {
	if m.next < len(m.mfs) {
		m.next += 1
		defer func() {
			if err := recover(); err != nil {
				elog.Fatal("%v", err)
				elog.Fatal("%v", string(Stack()))
			}
		}()
		m.mfs[m.next-1](m)

	}
	return m
}
