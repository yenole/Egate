package assist

import (
	"sync"
)

type EList struct {
	es  []interface{} //数据
	mtx sync.RWMutex  //互斥
}

func NewElist() *EList {
	return NewEListByLen(10)
}

func NewEListByLen(len int) *EList {
	l := new(EList)
	l.es = make([]interface{}, 0, len)
	return l
}

func (l *EList) Len() int {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	return len(l.es)
}

func (l *EList) Push(e ...interface{}) *EList {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.es = append(l.es, e...)
	return l
}

func (l *EList) Insert(i int, e ...interface{}) *EList {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.es = append(l.es[:i], append(e, l.es[i:]...)...)
	return l
}

func (l *EList) Remove(i int) interface{} {
	if i < 0 || i > l.Len()-1 {
		return nil
	}
	l.mtx.Lock()
	defer l.mtx.Unlock()
	r := l.es[i]
	l.es = append(l.es[:i], l.es[i+1:]...)
	return r
}

func (l *EList) Get(i int) interface{} {
	if i < 0 || i > l.Len()-1 {
		return nil
	}
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	return l.es[i]
}

func (l *EList) Pop() interface{} {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	r := l.es[l.Len()-1]
	l.es = l.es[:l.Len()-1]
	return r
}

func (l *EList) ForEach(f func(v interface{}, i int) bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	for i, v := range l.es {
		if !f(v, i) {
			break
		}
	}
}

func (l *EList) Lock() *EList {
	l.mtx.RLock()
	return l
}

func (l *EList) Unlock() *EList {
	l.mtx.RUnlock()
	return l
}

func (l *EList) Range() []interface{} {
	return l.es
}
