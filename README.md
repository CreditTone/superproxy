# 基于golang协程实现的超级代理，真正优雅使用代理池

### 效果展示
只要配置一个代理，发出去的ip都是动态切换的。爬虫不需考虑代理切换的问题。

```shell
stephendeMacBook-Pro:src stephen$ curl -x 127.0.0.1:12001 "http://ip.sb"
125.79.233.15
stephendeMacBook-Pro:src stephen$ curl -x 127.0.0.1:12001 "http://ip.sb"
60.182.32.64
stephendeMacBook-Pro:src stephen$ curl -x 127.0.0.1:12001 "http://ip.sb"
curl: (56) Recv failure: Connection reset by peer
stephendeMacBook-Pro:src stephen$ curl -x 127.0.0.1:12001 "http://ip.sb"
221.230.170.109
```

### 编译linux可执行文件
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build