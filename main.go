package main

import (
	"fmt"
	"github.com/dearkk/component/klog"
	"github.com/dearkk/component/plugins/apple"
	"os"
	"plugin"
)

func main() {
	p, err := plugin.Open("/Users/kun/workspace/bean/apple/plugin/apple.so")
	if err != nil {
		fmt.Println("error open plugin: \n", err)
		os.Exit(-1)
	}
	s, err := p.Lookup("Load")
	if err != nil {
		fmt.Println("error lookup Hello: \n", err)
		os.Exit(-1)
	}
	if load, ok := s.(func() interface{}); ok {
		fuc := load()
		a := fuc.(apple.Hello)
		klog.Printf("add: %d\n", a.Add(1, 5))
	}
}
