package day02

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 注册路由方法测试
func TestRouter_AddRoute(t *testing.T) {
	//第一步是构造路由树
	//第二步是验证路由树
	testRoutes := []struct{
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
		{
			method:http.MethodGet,
			path:"/order/detail/:id",
		},

	}
	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	// 下面注释的是测AddRoute添加方法能否成功，为了避免验证下面函数对重复路由注册验证的支持，测路由重复注册把下面注册路由先注释掉
	for _,route := range testRoutes {
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
					"order":&node{
						path:"order",
						children: map[string]*node{
							"detail":&node{
								path:"detail",
								paramChild: &node{
									path: ":id",
									handler: mockHandler,
								},
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
	assert.True(t,ok,msg)
	/*assert.Panics(t,func() {
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
	})*/
	r = newRouter()
	r.AddRoute(http.MethodGet,"/a/*",mockHandler)
	assert.Panicsf(t,func ()  {
		r.AddRoute(http.MethodGet,"/a/:id",mockHandler)
	},"web: 不允许同时注册路径参数和通配符匹配，已有通配符匹配")
	r = newRouter()
	r.AddRoute(http.MethodGet,"/a/:id",mockHandler)
	assert.Panicsf(t,func ()  {
		r.AddRoute(http.MethodGet,"/a/*",mockHandler)
	},"web: 不允许同时注册路径参数和通配符匹配，已有路径参数匹配")
}
func TestRouter_findRoute(t *testing.T) {
	testRoute := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/origin/home",
		},
		{
			method: http.MethodPost,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/login/:username",
		},
	}
	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}
	for _, route := range testRoute {
		r.AddRoute(route.method, route.path, mockHandler)
	}
	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		info  *matchInfo
	}{
		{
			name:      "user home",
			method:    http.MethodGet,
			path:      "/user/home",
			wantFound: true,
			info: &matchInfo{
				n:&node{
					handler: mockHandler,
					path:    "home",
				},
			},
		},
		{
			// username路径参数匹配
			name : "login username",
			method: http.MethodGet,
			path: "/login/dad",
			wantFound: true,
			info: &matchInfo{
				n:&node{
					path:":username",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"username":"dad",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t,tc.info.pathParams,info.pathParams)
			msg, ok := tc.info.n.equal(info.n)
			assert.True(t, ok, msg)

		})
	}
}
func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的http method"), false
		}
		// v,dst要比较相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不相同"), false
	}

	if n.starChid != nil {
		msg,ok := n.starChid.equal(y.starChid)
		if !ok{
			return msg,ok
		}
	}
	if n.paramChild != nil{
		msg,ok:=n.paramChild.equal(y.paramChild)
		if !ok{
			return msg,ok
		}
	}
	//handler比较
	nhandler := reflect.ValueOf(n.handler)
	yhandler := reflect.ValueOf(y.handler)
	if nhandler != yhandler {
		return fmt.Sprintf("handler不相等"), false
	}

	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点不存在"), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}
