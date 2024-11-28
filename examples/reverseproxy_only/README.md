# SRP 仅反向代理样例

本样例展示了仅进行反向代理，将本地服务通过其他方式暴露的样例。

首先通过以下命令执行该样例的 SRP 服务和本地样例服务：

```bash
# Shell 1
go run ./examples/reverseproxy_only/main.go

# Shell 2
go run ./examples/http/main.go
```

然后进行反向代理：

```bash
ssh -NR /www.example.com/80:127.0.0.1:8000 -p 8022 127.0.0.1
```

通过以下命令访问即可将请求反向代理到本地样例服务上：

```bash
curl -H 'Host: www.example.com' http://127.0.0.1:8008
```

还可以自定义更多的反向代理规则，例如根据 URL 进行反向代理等。
