package proto

import (
	"github.com/golang/protobuf/proto"
	"reflect"
)

type ProtoMsgParse struct {}

func (p *ProtoMsgParse) Marshal(m interface{}) ([]byte, error) {
	return proto.Marshal(m.(proto.Message))
}

func (p *ProtoMsgParse) Unmarshal(buf []byte, t reflect.Type) (interface{}, error) {
	m := reflect.New(t.Elem()).Interface()
	return m, proto.UnmarshalMerge(buf, m.(proto.Message))
}
