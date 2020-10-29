/*
Copyright © 2020 iiusky sky@03sec.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

// response修改
func ModifyResponse(resp *http.Response) error {
	var respBodyByte []byte

	if resp.Header.Get("Location") != "" {
		resp.Header.Set("Location", strings.ReplaceAll(resp.Header.Get("Location"),
			fmt.Sprintf("%s://%s", scheme, host), WebDomain))
	}

	respBodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Header.Del("Content-Encoding")
		respBodyByte = unGzip(respBodyByte)
	}

	// 实时打印代理情况
	fmt.Println(time.Now().String(), resp.Request.Method, resp.Request.URL.String())

	respBodyByte = bytes.Replace(respBodyByte, []byte("<body"), []byte(fmt.Sprintf(iframeContent, WebDomain, WebRtcPath)), -1)

	respbody := ioutil.NopCloser(bytes.NewReader(respBodyByte))
	resp.Body = respbody
	resp.ContentLength = int64(len(respBodyByte))
	resp.Header.Set("Content-Length", strconv.Itoa(len(respBodyByte)))

	return nil
}

// 反向代理
func handleReverseProxy(ctx *gin.Context) {

	director := func(req *http.Request) {
		req.URL.Scheme = scheme
		req.URL.Host = host
		req.Host = host
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ModifyResponse = ModifyResponse
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
