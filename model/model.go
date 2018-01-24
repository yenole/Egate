package model

import (
	"egate/gate"
	"net"
	"reflect"
)

var models = NewPool()

func init() {
	EPools.Push(-1, models)
}

type Agenter interface {
	GetAgent() gate.Agent
	SetAgent(agent gate.Agent)
}

type Model struct {
	gate.Agent
}

func (m *Model) GetAgent() gate.Agent {
	return m.Agent
}

func (m *Model) SetAgent(agent gate.Agent) {
	m.Agent = agent
}

func (m *Model) AgentWriteMsg(msg interface{}) {
	if m.Agent != nil {
		m.Agent.AgentWriteMsg(msg)
	}
}

func (m *Model) RemoteAddr() net.Addr {
	if m.Agent != nil {
		return m.Agent.RemoteAddr()
	}
	return nil
}

func (m *Model) ModelSetup(der gate.IDer, agent gate.Agent) {
	if p := models.Get(der.ID()); p != nil {
		m := p.(Agenter)
		ModelSwitchAgent(m.GetAgent(), agent)
		m.SetAgent(agent)
		callMethod(p, EV_AGENT_ATTEMPT)
	} else {
		m.SetAgent(agent)
		ModelAgent(agent, der)
		models.Push(der.ID(), der)
		callMethod(der, EV_AGENT_OPEN)
	}
}

func (m *Model) ModelUnSetup(der gate.IDer) {
	models.Pop(der.ID())
}

func (m *Model) Online() bool {
	return m.Agent != nil
}

func callMethod(m interface{}, method string) {
	if m != nil && len(method) > 0 {
		if fn := reflect.ValueOf(m).MethodByName(method); fn.Kind() == reflect.Func {
			fn.Call([]reflect.Value{})
		}
	}
}

func ModelByID(id uint64) gate.IDer {
	if p := models.Get(id); p != nil {
		return p.(gate.IDer)
	}
	return nil
}

func ModelSetupNum() int {
	return models.Len()
}
