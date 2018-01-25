package json

import (
	"egate/model"
	"encoding/json"
	"reflect"
)

type JsonMsgParse struct{}

func init() {
	model.MsgParseDiv(new(JsonMsgParse))
}

func (p *JsonMsgParse) Marshal(m interface{}) ([]byte, error) {
	return json.Marshal(m)
}

func (p *JsonMsgParse) Unmarshal(buf []byte, t reflect.Type) (interface{}, error) {
	m := reflect.New(t.Elem()).Interface()
	return m, json.Unmarshal(buf, m)
}
