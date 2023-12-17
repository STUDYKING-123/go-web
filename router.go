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
// findRoute:路由查找
func (r *router)findRoute(method string,path string)(*matchInfo, bool){
	root,ok := r.trees[method]
	if !ok {
		return nil,false
	}
	if path == "/"{
		return &matchInfo{
			n:root,
		},true
	}
	path = strings.Trim(path,"/")
	segs := strings.Split(path,"/")
	var pathParams map[string]string
	for _,seg := range segs{
		child,paramChild,found := root.childOf(seg)
		if !found {
			return nil,false
		}
		if paramChild{//命中路径参数
			if pathParams == nil{
				pathParams = make(map[string]string)
			}
			// path是 :id形式，我们去除前面的:
			pathParams[child.path[1:]] = seg
		}
		root = child
	}
	return &matchInfo{
		n:root,
		pathParams: pathParams,
	}, true
	//我们这种匹配考虑注册这样一个路由 /user/home，那么我们既能找到user节点也能找到home节点但是user节点上没有注册的handler
	//return root,root.handler != nil
}
// childOf优先考虑静态匹配，匹配不上考虑通配符匹配
// 第一个返回值是子节点， 第二个是标记返回子节点到底是否是路径参数 第三个标记是否命中
func (n *node) childOf(path string) (*node,bool,bool){
	if n.children == nil{
		if n.paramChild != nil{
			return n.paramChild,true,true
		}
		//没有子节点也要考虑通配符节点
		return n.starChid,false,n.starChid!=nil
	}
	child,ok := n.children[path]
	if !ok{
		if n.paramChild != nil{
			return n.paramChild,true,true
		}
		return n.starChid,false,n.starChid!=nil
	}
	return child,false,ok
}
func (n *node) childOrCreate(seg string) *node {
	if seg[0] == ':'{
		if n.starChid != nil{
			panic("web: 不允许同时注册路径参数和通配符匹配，已有通配符匹配")
		}
		n.paramChild = &node{
			path:seg,
		}
		return n.paramChild
	}

	if seg == "*" {
		if n.paramChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配，已有路径参数匹配")
		}
		n.starChid = &node{
			path:seg,
		}
		return n.starChid
	}
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

	//通配符匹配节点
	starChid *node

	// 加一个路径参数
	paramChild *node

	// 用户注册的业务逻辑
	handler HandleFunc
}

type matchInfo struct{
	n *node
	pathParams map[string]string
}
