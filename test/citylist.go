// citylist
package test

import (
	"hope/http/request"
	"regexp"
	"strings"
)

const cityListRe = `{linkContent:"([\s\S]*?)",linkURL:"([\s\S]*?)"}`

// {linkContent:"资阳",linkURL:"http://m.zhenai.com/zhenghun/ziyang1"},
func ParseCityList(contents []byte) *request.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1)
	result := &request.ParseResult{}
	// limit := 4
	for _, m := range matches {
		// result.Items = append(result.Items, "City "+string(m[1]))
		result.Requests = append(result.Requests, &request.Request{
			Url:       strings.Replace(string(m[2]), "//m.", "//www.", -1),
			ParseFunc: ParseCity,
		})
		// limit--
		// if limit == 0 {
		// 	break
		// }
	}
	return result
}
