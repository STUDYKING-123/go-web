package day02

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter_AddRoute(t *testing.T) {
	//第一步是构造路由树
	//第二步是验证路由树
	/*testRoutes := []struct{
		method string
		path string
	}{
		{
			method:http.MethodGet,
			path:"/",
		},
		{
			method:http.MethodGet,
			path:"/user",
		},
		{
			method:http.MethodGet,
			path:"/user/home",
		},
		{
			method:http.MethodGet,
			path:"/origin/home",
		},
		{
			method:http.MethodPost,
			path:"/user",
		},
		
	}*/
	var mockHandler HandleFunc = func(ctx *Context) {}  
	r := newRouter()
	// 下面注释的是测AddRoute添加方法能否成功，为了避免验证下面函数对重复路由注册验证的支持，测路由重复注册把下面注册路由先注释掉
	/*for _,route := range testRoutes {
		r.AddRoute(route.method,route.path,mockHandler)
	}

	//在这里断言路由树和你预期的一样,不用assert是因为HandleFunc是我们定义的方法，而方法在go里面是不可比的
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:"/",
				handler: mockHandler,
				children:map[string]*node{
					"user": &node{
						path:"user",
						handler: mockHandler,
						children: map[string]*node{
							"home":&node{
								path:"home",
								handler:mockHandler,
							},
						},
					},
					"origin":&node{
						path:"origin",
						children: map[string]*node{
							"home":&node{
								path:"home",
								handler: mockHandler,
							},
						},
					},
				},
			},
			http.MethodPost:&node{
				path:"/",
				children: map[string]*node{
					"user":&node{
						path:"user",
						handler: mockHandler,
					},
				},
			},
		},
	}
	msg,ok := wantRouter.equal(r)
	assert.True(t,ok,msg)*/
	assert.Panics(t,func() {
		r.AddRoute(http.MethodGet,"",mockHandler)
	})
	assert.Panics(t,func ()  {
		r.AddRoute(http.MethodGet,"sadjl",mockHandler)
	})
	assert.Panics(t,func() {
		r.AddRoute(http.MethodGet,"/asjdlkadj/",mockHandler)
	})
	assert.Panics(t,func() {
		r.AddRoute(http.MethodGet,"/us//jkl",mockHandler)
	},"web :中间不能有连续 /")
	r.AddRoute(http.MethodGet,"/",mockHandler)
	assert.PanicsWithValue(t,"web :路由冲突[/]",func() {
		r.AddRoute(http.MethodGet,"/",mockHandler)
	})
	r.AddRoute(http.MethodGet,"/as",mockHandler)
	assert.PanicsWithValue(t,"web :路由冲突[/as]",func() {
		r.AddRoute(http.MethodGet,"/as",mockHandler)
	})
}
func (r *router)equal(y *router) (string,bool) {
	for k,v := range r.trees{
		dst,ok := y.trees[k]
		if !ok{
			return fmt.Sprintf("找不到对应的http method"),false
		}
		// v,dst要比较相等
		msg,equal := v.equal(dst)
		if !equal{
			return msg,false
		}
	}
	return "",true
}

func (n *node)equal(y *node) (string,bool) {
	if n.path != y.path{
		return fmt.Sprintf("节点路径不匹配"),false
	}
	if len(n.children) != len(y.children){
		return fmt.Sprintf("子节点数量不相同"),false
	}

	//handler比较
	nhandler := reflect.ValueOf(n.handler)
	yhandler := reflect.ValueOf(y.handler)
	if nhandler != yhandler {
		return fmt.Sprintf("handler不相等"),false
	}

	for path,c := range n.children{
		dst ,ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点不存在"),false
		}
		msg,ok := c.equal(dst)
		if !ok {
			return msg,false
		}
	}
	return "", true
}