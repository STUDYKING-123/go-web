package day02

import (
	"fmt"
	"testing"
	"net/http"
)


func TestHTTPServer_ServeHTTP(t *testing.T){
	server := NewHttpServer()
	server.middles = [] Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个before")
				next(ctx)
				fmt.Println("第一个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个before")
				next(ctx)
				fmt.Println("第二个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三个中断")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第四个before")
				next(ctx)
				fmt.Println("第四个after")
			}
		},
	}
	server.ServeHTTP(nil,&http.Request{})
}