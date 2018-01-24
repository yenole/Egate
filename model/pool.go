package model

import "sync"

var EPools = NewPool()

type Pool struct {
	sync.RWMutex
	values map[interface{}]interface{}
}

func NewPool() *Pool {
	return &Pool{values: make(map[interface{}]interface{})}
}

func (p *Pool) Exist(key interface{}) bool {
	if key != nil {
		p.RLock()
		defer p.RUnlock()
		_, ok := p.values[key]
		return ok
	}
	return false
}

func (p *Pool) Get(key interface{}) interface{} {
	if key != nil {
		p.RLock()
		defer p.RUnlock()
		if v, ok := p.values[key]; ok {
			return v
		}
	}
	return nil
}

func (p *Pool) Push(key, value interface{}) interface{} {
	if key != nil {
		p.Lock()
		defer p.Unlock()
		if v, ok := p.values[key]; ok {
			p.values[key] = value
			return v
		}
		p.values[key] = value
	}
	return nil
}

func (p *Pool) Pop(key interface{}) interface{} {
	if key != nil {
		p.RLock()
		defer p.RUnlock()
		if v, ok := p.values[key]; ok {
			delete(p.values, key)
			return v
		}
	}
	return nil
}

func (p *Pool) Len() int {
	p.RLock()
	defer p.RUnlock()
	return len(p.values)
}
