//+build !dws

package egate

import (
	"egate/gate/network"
	"egate/gate"
	"errors"
)

// WS接受
func WsAccpet(egt *gate.Egate, addr string) {
	t := network.NewWsNetworkAccpet()
	if t.Init(make(network.Config).Address(addr).Build()) {
		egt.Accpet(t)
		return
	}
	panic(errors.New("ws accpet init error!"))
}
