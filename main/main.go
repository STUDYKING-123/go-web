package main

import (
	"day02"
	"net/http"
)

func main() {
	server := &day02.HttpServer{}
	server.AddRoute(http.MethodGet,"/",func(ctx *day02.Context) {})
	server.Start(":8080")
}
