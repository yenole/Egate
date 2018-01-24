package egate

import (
	"egate/elog"
	"egate/gate"
	"egate/gate/network"
	"egate/gate/packet"
	"egate/gate/route"
	"egate/model"
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const (
	IO_MODE_ROUTE uint = iota
	IO_MODE_MODEL
)

// 默认进出模式
func Mode(egt *gate.Egate, mode uint) {
	if egt != nil {
		egt.In(packet.In)
		if mode == IO_MODE_MODEL {
			egt.In(model.In)
			egt.Out(model.Out)
		}
		egt.In(route.In)
		egt.Out(packet.Out)
	}
}

// TCP接受
func TcpAccpet(egt *gate.Egate, addr string) {
	t := network.NewTcpNetworkAccpet()
	if t.Init(make(network.Config).Address(addr).Build()) {
		egt.Accpet(t)
		return
	}
	panic(errors.New("tcp accpet init error!"))
}

// 运行工作
func Work(egt *gate.Egate) {
	if egt != nil {
		var wg sync.WaitGroup
		egt.Work(&wg)
		elog.Info("egate %s running!", gate.Version())
		wg.Wait()
	}
}

// 守护进程
func Watch() {
	for i, v := range os.Args {
		if strings.EqualFold(v, "--watch") {
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			for {
				elog.Info("egate watch running!")
				exec.Command(os.Args[0], os.Args[1:]...).Run()
			}
		}
	}
}
