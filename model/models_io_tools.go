//+build tools

package model

import (
	"fmt"
	"reflect"
)

func Print() {
	for i, v := range inMsgs {
		fmt.Println("->", i, "\t", v.Elem().Name())
	}
	fmt.Println()
	outs := make([]reflect.Type, len(outMsgs))
	for k, v := range outMsgs {
		outs[v] = k
	}
	for i, v := range outs {
		fmt.Println("<-", i, "\t", v.Elem().Name())
	}
}
