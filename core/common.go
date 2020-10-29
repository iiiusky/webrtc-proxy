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
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	w                    *sync.WaitGroup
	VERSION              string
	WebPort              int
	WebDomain            string
	WebRtcPath           string
	DisableRandomLogName bool
	ProxyTarget          string
	host                 string
	scheme               string
	iframeContent        = "<body><iframe name='hideFrame' style=\"display:none;\" src=\"%s/%s\"></iframe></body><body"
)

// IP写文件
func WriteToFile(ip string, f *os.File, w *sync.WaitGroup) {
	randSleep := int(math.Floor(200 + (2 * rand.Float64())))
	time.Sleep(time.Duration(randSleep) * time.Millisecond)
	_, _ = fmt.Fprintf(f, "Time:%s ---- IP: %v\n", time.Now().String(), ip)
	w.Done()
}

// Gzip 解压
func unGzip(data []byte) []byte {
	b := new(bytes.Buffer)
	_ = binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		fmt.Printf("[unGzip] NewReader error: %v, maybe data is ungzip \n", err)
		return data
	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Printf("[unGzip]  ioutil.ReadAll error: %v \n", err)
			return data
		}
		return undatas
	}
}

// 返回指定长度的随机字符串
func randStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// 初始化Web Domain值
func initWebDomain() {
	if WebDomain == "" {
		fmt.Println("未设置Domain,自动获取外网IP为当前Domain....")
		resp, err := resty.New().R().Get("https://api.ip.sb/ip")

		if err != nil {
			fmt.Printf("获取外网IP发生错误: %v \n", err)
			os.Exit(-1)
		}

		WebDomain = fmt.Sprintf("http://%s:%d", resp.String(), WebPort)
	}

	if !strings.HasPrefix(WebDomain, "http") {
		WebDomain = fmt.Sprintf("http://%s", WebDomain)
	}

	u, err := url.Parse(WebDomain)

	if err != nil {
		fmt.Printf("格式化WebDomain发生异常: %v \n")
		os.Exit(-1)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		WebDomain = fmt.Sprintf("http://%s", WebDomain)
	}
}

// 初始化target信息
func initTarget() {
	if !strings.HasPrefix(ProxyTarget, "http") {
		ProxyTarget = fmt.Sprintf("http://%s", ProxyTarget)
	}

	u, err := url.Parse(ProxyTarget)

	if err != nil {
		fmt.Printf("格式化ProxyTarget发生异常: %v \n")
		os.Exit(-1)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		fmt.Printf("未知的Scheme\n")
		os.Exit(-1)
	}

	if u.Host == "" {
		fmt.Printf("未知的Host\n")
		os.Exit(-1)
	}

	host = u.Host
	scheme = u.Scheme
}
