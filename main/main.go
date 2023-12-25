package main

import (
	"day02"
	"net/http"
)
func hello(ctx *day02.Context){
	ctx.RespJOSN(200,"你好世界")
}
func main() {
	server := day02.NewHttpServer()
	//server := &day02.HttpServer{}
	server.AddRoute(http.MethodGet,"/",hello)
	server.AddRoute(http.MethodPost,"/user",hello)
	server.Start(":8080")
}
