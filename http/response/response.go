// response
package response

import (
	"hope/http/request"
)

type Response struct {
	Status     string
	StatusCode int
	Proto      string
	Headers    map[string]string
	Body       []byte
	Cookies    []string
	ParseFunc  func(contents []byte) *request.ParseResult
}
