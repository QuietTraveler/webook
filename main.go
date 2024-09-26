package main

import (
	"basic-go/webook/internal/web"
)

func main() {
	server := web.RegisterRoutes()
	err := server.Run(":8080")
	if err != nil {
		panic(err)
		return
	}
}
