// test
package test

import (
	"hope/http/request"
	"regexp"
)

var urlRe = regexp.MustCompile(`<div[\s\S]*?<a target="_blank" title="[\s\S]*?" href="([\s\S]*?)">([\s\S]*?)</a>[\s\S]*?</div>`)

func ParseUrl(contents []byte) *request.ParseResult {
	result := &request.ParseResult{}
	matches := urlRe.FindAllSubmatch(contents, -1)
	for _, match := range matches {
		result.Items = append(result.Items, match[2])
		result.Requests = append(result.Requests, &request.Request{
			Url:       "https://www.baidu.com" + string(match[1]),
			Headers:   map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36"},
			ParseFunc: ParseUrl,
		})
	}
	return result
}
