//+build client

package egate

import (
	"egate/gate"
	"egate/gate/packet"
	"egate/model"
	"egate/gate/route"
)

// 客户端默认进出模式
func ClientMode(clt *gate.Client, mode uint) {
	if clt != nil {
		clt.In(packet.In)
		if mode == IO_MODE_MODEL {
			clt.In(model.In)
			clt.Out(model.Out)
		}
		clt.In(route.In)
		clt.Out(packet.Out)
	}
}