package model

import (
	"egate/gate"
	"reflect"
	"sync"
)

const (
	EV_AGENT_OPEN    = "OnAgentOpen"
	EV_AGENT_CLOSE   = "OnAgentClose"
	EV_AGENT_ATTEMPT = "OnAgentAttempt"
)

var (
	amtx   sync.RWMutex
	agents = map[gate.Agent]interface{}{}
)

func ModelAgent(agent gate.Agent, model interface{}) bool {
	amtx.Lock()
	defer amtx.Unlock()
	if _, ok := agents[agent]; !ok {
		agents[agent] = model
		return true
	}
	return false
}

func ModelUnAgent(agent gate.Agent) bool {
	amtx.Lock()
	defer amtx.Unlock()
	if m, ok := agents[agent]; ok {
		callMethod(m, EV_AGENT_CLOSE)
		m.(Agenter).SetAgent(nil)
		delete(agents, agent)
		return true
	}
	return false
}

func ModelSwitchAgent(agent gate.Agent, agent2 gate.Agent) {
	amtx.Lock()
	defer amtx.Unlock()
	if m, ok := agents[agent]; ok {
		delete(agents, agent)
		agents[agent2] = m
		agent.Close()
	}
}

func ModelMethod(agent gate.Agent, idx int) reflect.Value {
	amtx.RLock()
	defer amtx.RUnlock()
	if m, ok := agents[agent]; ok {
		return reflect.ValueOf(m).Method(idx)
	}
	return reflect.Value{}
}
