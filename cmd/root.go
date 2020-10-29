/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"webrtc/core"
)

var rootCmd = &cobra.Command{
	Use:   "webrtc",
	Short: "反向代理目标网站+webrtc 获取真实IP",
	Long:  `反向代理目标网站+webrtc 获取真实IP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if core.ProxyTarget == "" {
			return errors.New("需要代理的目标网站不允许为空")
		}
		core.WebServer()
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&core.ProxyTarget, "target", "t", "", "需要代理的目标网站,,比如https://www.baidu.com")
	rootCmd.Flags().IntVarP(&core.WebPort, "port", "p", 80, "服务开启的端口,默认开启80端口")
	rootCmd.Flags().StringVarP(&core.WebDomain, "domain", "d", "",
		"服务绑定的域名(IP),如不指定则自动获取外网IP,目标网站代理后iframe的webrtc地址以此为准")
	rootCmd.Flags().StringVar(&core.WebRtcPath, "path", "webrtc", "webrtc访问的路由,默认为webrtc")
	rootCmd.Flags().BoolVar(&core.DisableRandomLogName, "disable-random-log", false,
		"是否关闭日志存放文件名随机生成,默认为false,如果为true则为[webrtc.log]、[access.log]创建,"+
			"否则为[随机字符串+webrtc.log]、[随机字符串+webrtc.log]生成.")
}
