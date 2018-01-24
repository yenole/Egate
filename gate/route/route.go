package route

import (
	"egate/gate"
	"reflect"
)

var (
	MsgHandles = map[uint32]reflect.Value{}
)

func Handle(id uint32, fun interface{}) bool {
	if fun != nil && reflect.ValueOf(fun).Kind() == reflect.Func {
		if _, ok := MsgHandles[id]; !ok {
			MsgHandles[id] = reflect.ValueOf(fun)
			return true
		}
	}
	return false
}

func In(m *gate.Middleware) {
	if typ, ok := MsgHandles[m.Msg.Id]; !m.IsAbort() && ok {
		args := []reflect.Value{
			reflect.ValueOf(m.Msg.Msg),
		}
		switch typ.Type().NumIn() {
		case 2:
			args = append(args, reflect.ValueOf(m.Agent))
		}
		typ.Call(args)
	}
	m.Next()
}
