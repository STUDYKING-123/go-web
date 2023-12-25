package day02

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)
type Context struct{
	Req *http.Request
	Response http.ResponseWriter
	pathParams map[string]string

	//命中路由，消息记录时用到，记录命中的路由
	MatchRoute string
	queryvalues url.Values
}
func (c *Context) BindJSON(val any) error {
	if val == nil {
		return errors.New("text:输入不能为空")

	}
	if c.Req.Body == nil{
		return errors.New("body 不能为空")
	}
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}
func (c *Context) FormValue(key string)(string,error){
	err := c.Req.ParseForm()
	if err != nil{
		return "",err
	}
	return c.Req.FormValue(key),nil
}
func (c *Context) QueryValue(key string)(string,error){
	if c.queryvalues == nil{
		c.queryvalues = c.Req.URL.Query()
	}
	vals ,ok:=c.queryvalues[key]
	if !ok{
		return "",errors.New("传入key 不存在")
	}
	// 用户无法知道是真的有值，但是值恰好是空字符串，还是没有值
	return vals[0],nil
}

func (c *Context) PathValue(key string) (string,error){
	val,ok := c.pathParams[key]
	if !ok {
		return "",errors.New("web:key不存在")
	}
	return val,nil
}

// 输出处理

func (c *Context) RespJOSN(status int,val any) error {
	data,err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Response.WriteHeader(status)
	//c.Response.Header().Set("Content-Type","application/json")
	n, err := c.Response.Write(data)
	if n != len(data) {
		return errors.New("web: 未写入全部数据")
	}
	return nil
}