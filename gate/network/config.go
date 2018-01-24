package network

import "egate/gate"

const (
	C_ADDRESS   = iota // 地址：127.0.0.1:8081
	C_PACKET           // 封包解包对象
	OPT_ADDRESS = "address"
)

type Config map[int]interface{}

func (t Config) Build() gate.Config {
	return gate.Config(t)
}

func (t Config) address() string {
	if addr, ok := t[C_ADDRESS]; ok {
		return addr.(string)
	}
	return ""
}

func (t Config) Address(addr string) Config {
	t[C_ADDRESS] = addr
	return t
}
