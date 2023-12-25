package day02

import (
	_ "fmt"
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
	middles []Middleware
}
type HTTPServerOption func(* HttpServer)
func NewHttpServer(opts ...HTTPServerOption)*HttpServer{
	res := &HttpServer{
		router: newRouter(),
	}
	for _,opt := range opts{
		opt(res)
	}
	return res
}
func ServerWithMiddleware (mdls ...Middleware) HTTPServerOption{
	return func(server *HttpServer){
		server.middles = mdls
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
	// 构造访问前需要访问的函数调用链，最后一个需要调用的函数是我们的serve,从后往前构造链，从前往后执行链
	root := h.serve
	for i := len(h.middles)-1 ;i >=0 ;i--{
		root = h.middles[i](root)
	}
	// 接下来就是查找路由并执行命中的业务逻辑
	root(ctx)
}
func (h *HttpServer)serve(ctx *Context){
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil{
		// 路由没有命中，返回404
		ctx.Response.WriteHeader(404)
		ctx.Response.Write([]byte("NOT FOUND"))
		return
	}
	ctx.pathParams = info.pathParams
	ctx.MatchRoute = info.n.route
	info.n.handler(ctx)
}
