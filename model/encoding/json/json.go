package json

import (
	"encoding/json"
	"reflect"
)

type JsonMsgParse struct{}

func (p *JsonMsgParse) Marshal(m interface{}) ([]byte, error) {
	return json.Marshal(m)
}

func (p *JsonMsgParse) Unmarshal(buf []byte, t reflect.Type) (interface{}, error) {
	m := reflect.New(t.Elem()).Interface()
	return m, json.Unmarshal(buf, m)
}
