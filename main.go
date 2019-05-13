// hope project main.go
package main

import (
	"hope/engine"
	"hope/http/request"
	"hope/test"
)

//https://httpbin.org/get
func main() {
	//proxy := "http://61.128.208.94:3128"
	e := engine.NewSpider("mysql", "root:@tcp(127.0.0.1:3306)/zhenai?charset=utf8", 5)
	e.Run(&request.Request{
		Url: "http://www.zhenai.com/zhenghun",
		// Proxy:     "http://61.128.208.94:3128",
		Headers:   map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36"},
		ParseFunc: test.ParseCityList,
	})
}
