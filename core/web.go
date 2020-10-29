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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jxskiss/ginregex"
	"os"
	"sync"
	"time"
	"webrtc/asset"
)

// Web服务
func WebServer() {
	r := gin.Default()
	w = new(sync.WaitGroup)
	webRtcLogPath := "webrtc.log"
	webLogPath := "access.log"

	if !DisableRandomLogName {
		webRtcLogPath = fmt.Sprintf("%s-%s", randStringRunes(6), webRtcLogPath)
		webLogPath = fmt.Sprintf("%s-%s", randStringRunes(6), webLogPath)
	}

	initWebDomain()
	initTarget()

	webRtcLogFile, err := os.OpenFile(webRtcLogPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("webRtcLogFile 文件打开失败", err)
		os.Exit(-1)
	}

	file, _ := os.Create(webLogPath)
	c := gin.LoggerConfig{
		Output:    file,
		SkipPaths: []string{""},
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				params.ClientIP,
				params.TimeStamp.Format(time.RFC1123),
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
			)
		},
	}

	r.Use(gin.LoggerWithConfig(c))

	r.POST("/test.json", func(c *gin.Context) {
		ip := c.PostForm("ips")
		fmt.Printf("time: %s  ---- ip: %s", time.Now().String(), ip)
		w.Add(1)
		WriteToFile(ip, webRtcLogFile, w)
		c.JSON(200, "hi~")
	})

	r.GET(fmt.Sprintf("/%s", WebRtcPath), func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		indexHtml, _ := asset.Asset("static/index.html")
		_, _ = c.Writer.Write(indexHtml)
		c.Writer.Header().Add("Accept", "text/html")
		c.Writer.Flush()
	})

	regexRouter := ginregex.New(r, nil)
	regexRouter.Any("^/.*$", handleReverseProxy)

	r.Run(fmt.Sprintf("0.0.0.0:%d", WebPort))
}
