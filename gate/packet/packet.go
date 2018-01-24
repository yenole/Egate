package packet

import (
	"bytes"
	"egate/elog"
	"egate/gate"
	"encoding/binary"
	"io"
)

func In(m *gate.Middleware) {
	lenBuf := make([]byte, 2)
	if _, err := m.Agent.Read(lenBuf); err != nil {
		if err != io.EOF {
			elog.Fatal("Agent.Read:%v", err)
		}
		m.Abort()
		m.Next()
		return
	}
	pLen := uint32(binary.LittleEndian.Uint16(lenBuf))
	msgBuf := make([]byte, pLen)
	if _, err := io.ReadFull(m.Agent, msgBuf); err != nil {
		elog.Fatal("io.ReadFull:%v", err)
		m.Abort()
		m.Next()
		return
	}
	m.Msg.Id = uint32(binary.BigEndian.Uint16(msgBuf[:2]))
	m.Msg.Msg = msgBuf[2:]
	m.Next()
}

func Out(m *gate.Middleware) {
	switch m.Msg.Msg.(type) {
	case []byte:
		msg := m.Msg.Msg.([]byte)
		lenBuf := make([]byte, 2)
		idBuf := make([]byte, 2)
		binary.BigEndian.PutUint16(idBuf, uint16(m.Msg.Id))
		binary.LittleEndian.PutUint16(lenBuf, uint16(len(msg)+2))
		buffer := bytes.NewBuffer(lenBuf)
		buffer.Write(idBuf)
		buffer.Write(msg)
		m.Agent.Write(buffer.Bytes())
	}
}
