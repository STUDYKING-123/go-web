package day02

import (
	"fmt"
	"strings"
)

type router struct {
	// Beego Gin HTTP method对应一棵树

	// http method对应一个树的根节点
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}
// path必须以 / 开头，不能以 / 结尾， 中间不能有连续的 //
func (r *router) AddRoute(method string, path string, handlefunc HandleFunc) {
	if path ==""{
		panic("web: 路径不能为空字符串")
	}
	// 开头结尾校验
	if path[0] != '/' {
		panic("web :路径必须以 / 开头")
	}
	if path != "/" && path[len(path)-1] == '/'{
		panic("web :不能以 / 结尾")
	} 
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	if path == "/"{
		if root.handler != nil {
			panic("web :路由冲突[/]")
		}
		root.handler = handlefunc
		return
	}

	// 切割这个路径
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		if seg == ""{
			panic("web :中间不能有连续 /")
		}
		//递归下去，找准位置
		children := root.childOrCreate(seg)
		root = children
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web :路由冲突[%s]",path))
	}
	root.handler = handlefunc
}

func (n *node) childOrCreate(seg string) *node {
	if n.children == nil {
		n.children = make(map[string]*node)
		res := &node{
			path: seg,
		}
		n.children[seg] = res
		return res
	}
	res, ok := n.children[seg]
	if !ok {
		//添加路由时没有当前节点要新建当前节点
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

type node struct {
	path string
	// 子path到子节点的映射
	children map[string]*node

	// 用户注册的业务逻辑
	handler HandleFunc
}
