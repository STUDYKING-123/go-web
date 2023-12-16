package day02
import "net/http"
type Context struct{
	Req *http.Request
	Response http.ResponseWriter
}