package accesslog

import (
	"day02"
	"fmt"
	"net/http"
	"testing"
)


func TestMiddlewareBuilder(t *testing.T){
	builder := MiddlewareBuilder{}
	mdls := builder.LogFunc(func (log string)  {
		fmt.Println(log)
	}).Build()
	server := day02.NewHttpServer(day02.ServerWithMiddleware(mdls))
	server.POST("/a/b/*",func(ctx *day02.Context) {
		fmt.Println("hello world")
	})
	req, _ := http.NewRequest(http.MethodPost,"/a/b/c",nil)
	req.Host = "localhost"
	server.ServeHTTP(nil,req)
}