// downloader
package download

import (
	"bufio"
	"fmt"
	"hope/http/request"
	"hope/http/response"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/text/encoding"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func Download(r *request.Request) (*response.Response, error) {
	//创建客户端
	var transport *http.Transport
	var client *http.Client
	// 添加proxy
	// 设置超时时间
	if r.Proxy != "" {
		urli := url.URL{}
		urlProxy, _ := urli.Parse(r.Proxy)
		transport = &http.Transport{
			Proxy: http.ProxyURL(urlProxy),
		}
		client = &http.Client{
			Transport: transport,
			Timeout:   time.Duration(r.Timeout),
		}
	} else {
		client = &http.Client{Timeout: time.Duration(r.Timeout)}
	}
	// 请求方法，url，主体，生成Request
	req, err := http.NewRequest(r.Method, r.Url, r.Body)
	if err != nil {
		return &response.Response{}, fmt.Errorf("<Hope downloader>: create request fail <%s>", err)
	}
	// 添加Headers
	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}
	// 发起请求，获取Response
	resp, err := client.Do(req)
	if err != nil {
		return &response.Response{}, fmt.Errorf("<Hope downloader>: request fail <%s>", err)
	}
	if resp.StatusCode != http.StatusOK {
		return &response.Response{}, fmt.Errorf("<Hope downloader>: wrong status code: %d", resp.StatusCode)
	}
	newResponse := &response.Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Proto:      resp.Proto,
	}
	//网页内容转码位utf8
	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	body, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		return &response.Response{}, fmt.Errorf("<Hope downloader>: fail read from response body %s", err)
	}
	newResponse.Body = body
	//cookie
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		newResponse.Cookies = append(newResponse.Cookies, cookie.String())
	}
	newResponse.ParseFunc = r.ParseFunc
	defer resp.Body.Close()
	return newResponse, nil
}

//自动发现网页编码
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("<Hope downloader>: transform encoding error %v", err)
		return unicode.UTF8
	}
	//通过Content-Type和前1024bit内容的编码来发现网页的编码
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
