package main

import (
	"ChildrenMath/pkg/settings"
	"ChildrenMath/router"
)

func main() {
	r := router.InitRouter()

	//err := r.Run(settings.ServerIp + ":" + settings.ServerPort)
	err := r.RunTLS(settings.ServerIp+":"+settings.ServerPort, "./cert/server.crt", "./cert/server.key")
	if err != nil {
		panic(err)
	}
}
