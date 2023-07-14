package main

import (
	"ChildrenMath/pkg/settings"
	"ChildrenMath/router"
)

func main() {
	r := router.InitRouter()

	err := r.Run(":" + settings.HttpPort)
	if err != nil {
		panic(err)
	}
}
