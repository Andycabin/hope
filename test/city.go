// city
package test

import (
	"hope/http/request"
	"regexp"
	"strings"
)

var cityRe = regexp.MustCompile(`<tr><th><a href="([\s\S]*?)" target="_blank">([\s\S]*?)</a></th></tr>`)
var nextUrlRe = regexp.MustCompile(`<li class="paging-item"><a href="(.*?)">下一页</a>`)

func ParseCity(contents []byte) *request.ParseResult {
	result := &request.ParseResult{}
	matches := cityRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, &request.Request{
			Url:       string(m[1]),
			Headers:   map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36"},
			ParseFunc: ParseProfile,
		})
	}
	next_matche := nextUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range next_matche {
		nexturl := strings.Split(string(m[1]), "href=\"")
		result.Requests = append(result.Requests, &request.Request{
			Url:       nexturl[len(nexturl)-1],
			Headers:   map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36"},
			ParseFunc: ParseCity,
		})
	}
	return result
}
