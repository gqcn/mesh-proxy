package main

import (
	_ "mesh-proxy/boot"

	"github.com/gogf/gf/os/gcmd"
	"mesh-proxy/app/service/command"
)

func main() {
	err := gcmd.BindHandleMap(map[string]func(){
		"proxy":  command.Proxy,
		"health": command.Health,
		"server": command.Server,
	})
	if err != nil {
		panic(err)
	}
	err = gcmd.AutoRun()
	if err != nil {
		panic(err)
	}
}
