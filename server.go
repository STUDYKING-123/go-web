package day02

import (
	"fmt"
	"net/http"
)
type HandleFunc func(ctx *Context)
type Server interface{
	Start(address string ) error
	http.Handler

	// AddRoute路由注册功能
	// method:http方法，path:URL路径,HandleFunc:业务处理逻辑
	AddRoute(method string,path string,handlefunc HandleFunc)
	// 不提供以下这种设计方案
	/*
	AddRoute(method string,path string,handlefunc... HandleFunc)
	实际上用户调用时
	AddRoute(method,path,func(ctx *Context){
		handle1(ctx)
		handle2(ctx)
	})
	与AddRoute(method,path,handle1,handle2)等效，我们将处理权交给用户，但是后面这种做法还可以允许用户不传入对应处理逻辑
	
	*/
}
type HttpServer struct{
	*router
}

func NewHttpServer()*HttpServer{
	return &HttpServer{
		router: newRouter(),
	}
}

var _ Server = &HttpServer{}
/*func (h *HttpServer)AddRoute(method string,path string,handlefunc HandleFunc){
	//注册对应方法+PATH到路由树
}*/
func (h *HttpServer)Get(path string,handlefunc HandleFunc){
	h.AddRoute(http.MethodGet,path,handlefunc)
}
func (h *HttpServer)POST(path string,handlefunc HandleFunc){
	h.AddRoute(http.MethodPost,path,handlefunc)
}
func (h *HttpServer)DEL(path string,handlefunc HandleFunc){
	h.AddRoute(http.MethodDelete,path,handlefunc)
}
func (h *HttpServer)PUT(path string,handlefunc HandleFunc){
	h.AddRoute(http.MethodPut,path,handlefunc)
}
func (h *HttpServer)Start(address string) error {
	return http.ListenAndServe(address,h)
}
// ServeHTTP是作为http包和web框架的关联点，需要在ServeHTTP包内部，
//执行：构建起Web框架的上下文，查找路由树，并执行命中路由的代码
func (h *HttpServer)ServeHTTP(w http.ResponseWriter,r *http.Request){
	ctx := &Context{
		Req: r,
		Response: w,
	}
	// 接下来就是查找路由并执行命中的业务逻辑
	h.serve(ctx)
}
func (h *HttpServer)serve(ctx *Context){
	fmt.Fprintf(ctx.Response,"hello")
}
