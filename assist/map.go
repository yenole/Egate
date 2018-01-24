package assist

import (
	"sync"
)

type EMap struct {
	es  map[interface{}]interface{} //数据
	mtx sync.RWMutex                //互斥
}

func NewEMap() *EMap {
	m := new(EMap)
	m.es = make(map[interface{}]interface{})
	return m
}

func (m *EMap) Len() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return len(m.es)
}

func (m *EMap) Push(key, value interface{}) *EMap {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.es[key] = value
	return m
}

func (m *EMap) Get(key interface{}) interface{} {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	if v, ok := m.es[key]; ok {
		return v
	}
	return nil
}

func (m *EMap) Pop(key interface{}) interface{} {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if v, ok := m.es[key]; ok {
		delete(m.es, key)
		return v
	}
	return nil
}

func (m *EMap) Keys() *EList {
	l := NewEListByLen(m.Len())
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	for k, _ := range m.es {
		l.Push(k)
	}
	return l
}

func (m *EMap) Valuse() *EList {
	l := NewEListByLen(m.Len())
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	for _, v := range m.es {
		l.Push(v)
	}
	return l
}

func (m *EMap) ForEach(f func(v, k interface{}) bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	for i, v := range m.es {
		if !f(v, i) {
			break
		}
	}
}

func (m *EMap) Lock() *EMap {
	m.mtx.RLock()
	return m
}

func (m *EMap) Unlock() *EMap {
	m.mtx.RUnlock()
	return m
}

func (m *EMap) Range() map[interface{}]interface{} {
	return m.es
}
