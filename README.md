# 反向代理+webrtc 神不知鬼不觉的获取真实IP

## Usage
```
Usage:
  webrtc [flags]
  webrtc [command]

Available Commands:
  help        Help about any command
  version     显示版本信息

Flags:
      --disable-random-log   是否关闭日志存放文件名随机生成,默认为false,如果为true则为[webrtc.log]、[access.log]创建,否则为[随机字符串+webrtc.log]、[随机字符串+webrtc.log]生成.
  -d, --domain string        服务绑定的域名(IP),如不指定则自动获取外网IP,目标网站代理后iframe的webrtc地址以此为准
  -h, --help                 help for webrtc
      --path string          webrtc访问的路由,默认为webrtc (default "webrtc")
  -p, --port int             服务开启的端口 (default 80)
  -t, --target string        需要代理的目标网站,比如https://www.baidu.com

Use "webrtc [command] --help" for more information about a command.
```
## 实现原理
通过反向代理目标网站,然后替换目标网站中返回的response包中`<body`的内容,使用iframe嵌套进去原始网页,从而神不知鬼不觉的获取到真实IP
