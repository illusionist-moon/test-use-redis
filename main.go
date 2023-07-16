package main

import (
	"ChildrenMath/pkg/settings"
	"ChildrenMath/router"
	"net/http"
)

func main() {
	r := router.InitRouter()

	//err := r.Run(settings.ServerIp + ":" + settings.ServerPort)

	//err := r.RunTLS(settings.ServerIp+":"+settings.ServerPort, "./cert/server.crt", "./cert/server.key")

	//server := &http.Server{
	//	Addr:    settings.ServerIp + ":" + settings.ServerPort,
	//	Handler: r,
	//}
	//err := server.ListenAndServeTLS("./cert/server.crt", "./cert/server.key")

	err := http.ListenAndServeTLS(settings.ServerIp+":"+settings.ServerPort,
		"./cert/server.crt",
		"./cert/server.key", r)

	if err != nil {
		panic(err)
	}
}
