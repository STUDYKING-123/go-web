package accesslog

import (
	"day02"
	"encoding/json"
)

type MiddlewareBuilder struct{
	logFunc func(log string)
}
func (m *MiddlewareBuilder) LogFunc(fn func(log string))*MiddlewareBuilder{
	m.logFunc = fn
	return m 
}

func (m MiddlewareBuilder) Build() day02.Middleware{
	return func(next day02.HandleFunc) day02.HandleFunc {
		return func(ctx *day02.Context) {
			// 要记录请求
			// 在defer里面才最终输出日志，因为:确保即便next里面发生了panic，也能将请求正确记录下来
			// 获取MatchRoute必须在执行了next之后才能获得，因为依赖于最终路由树匹配(HTTPServer.serve)
			defer func() {
				l := accessLog{
					Host: ctx.Req.Host,
					Route: ctx.MatchRoute,
					HTTPMethod: ctx.Req.Method,
					Path: ctx.Req.URL.Path,
				}
				data,_ := json.Marshal(l)
				m.logFunc(string(data))
			}()
			next(ctx)
		}
	}
}

type accessLog struct{
	Host string `json:"host,omitempty"`
	Route string `json: "route,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`
	Path string `json:"path,omitempty"`
}