//+build !dws

package network

import (
	"bytes"
	"github.com/gorilla/websocket"
	"io"
)

type WsConn struct {
	*websocket.Conn
	buffer *bytes.Buffer
}

func NewWsConn(conn *websocket.Conn) *WsConn {
	return &WsConn{Conn: conn, buffer: bytes.NewBuffer([]byte{})}
}

func (w *WsConn) Read(p []byte) (n int, err error) {
	if w.buffer == nil {
		w.buffer = bytes.NewBuffer([]byte{})
	}
	if w.buffer.Len() < len(p) {
		mt, b, err := w.ReadMessage()
		if err != nil {
			return 0, err
		} else if mt == websocket.CloseMessage {
			// 关闭消息
			return 0, io.EOF
		}
		w.buffer.Write(b)
	}
	return w.buffer.Read(p)
}

func (w *WsConn) Write(p []byte) (n int, err error) {
	return len(p), w.WriteMessage(websocket.BinaryMessage, p)
}
