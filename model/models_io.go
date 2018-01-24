package model

import (
	"egate/elog"
	"egate/gate"
	"egate/model/encoding/json"
	"egate/model/encoding/proto"
	"reflect"
)

const (
	MSG_PARSER_PROTO = iota //PROTOCBUF
	MSG_PARSER_JSON
)

var (
	msgParse  MsgParser
	inMsgs    = []reflect.Type{}
	outMsgs   = map[reflect.Type]uint32{}
	modelMsgs = map[uint32]int{}
)

type MsgParser interface {
	Marshal(m interface{}) ([]byte, error)
	Unmarshal(p []byte, t reflect.Type) (interface{}, error)
}

func MsgParse(id uint32) {
	switch id {
	case MSG_PARSER_PROTO:
		msgParse = &proto.ProtoMsgParse{}
	case MSG_PARSER_JSON:
		msgParse = &json.JsonMsgParse{}
	}
}

func MsgParseDiv(p MsgParser) {
	if p != nil {
		msgParse = p
	}
}

// 消息清空
func MsgClean() {
	inMsgs = []reflect.Type{}
	outMsgs = map[reflect.Type]uint32{}
	modelMsgs = map[uint32]int{}
}

//	模块消息进
func MsgModelIn(m interface{}, fs ...string) {
	if m != nil {
		t := reflect.ValueOf(m).Type()
		for _, v := range fs {
			if m, ok := t.MethodByName(v); ok {
				modelMsgs[uint32(len(inMsgs))] = m.Index
				inMsgs = append(inMsgs, m.Type.In(1))
			} else {
				elog.Fatal("%s %s method not found!", t.Elem().Name(), v)
			}
		}
	}
}

// 消息进
func MsgIn(msgs ...interface{}) {
	for _, v := range msgs {
		inMsgs = append(inMsgs, reflect.ValueOf(v).Type())
	}
}

//　消息出
func MsgOut(msgs ...interface{}) {
	for _, v := range msgs {
		outMsgs[reflect.ValueOf(v).Type()] = uint32(len(outMsgs))
	}
}

func In(m *gate.Middleware) {
	if m.IsAbort() {
		ModelUnAgent(m.Agent)
	} else if msgParse != nil && m.Msg.Msg != nil {
		msg, err := msgParse.Unmarshal(m.Msg.Msg.([]byte), inMsgs[m.Msg.Id])
		if err != nil {
			elog.Fatal("%v", err)
		} else if idx, ok := modelMsgs[m.Msg.Id]; ok {
			m.Msg.Msg = msg
			if tv := ModelMethod(m.Agent, idx); tv.Kind() == reflect.Func {
				tv.Call([]reflect.Value{reflect.ValueOf(m.Msg.Msg)})
				return
			}
		} else {
			m.Msg.Msg = msg
		}
	}
	m.Next()
}

func Out(m *gate.Middleware) {
	if msgParse != nil && m.Msg.Msg != nil {
		if id, ok := outMsgs[reflect.ValueOf(m.Msg.Msg).Type()]; ok {
			buf, err := msgParse.Marshal(m.Msg.Msg)
			if err != nil {
				elog.Fatal("%v", err)
			} else {
				m.Msg.Msg = buf
				m.Msg.Id = id
			}
		}
	}
	m.Next()
}
