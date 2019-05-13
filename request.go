// request
// Request数据结构
// 1.包含请求的必须信息
// 2.发起GET请求
// 3.返回响应
package request

import (
	"io"
)

type Request struct {
	Url       string
	Method    string
	Headers   map[string]string
	Body      io.Reader
	Proxy     string
	Timeout   int
	ParseFunc func([]byte) *ParseResult
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}
