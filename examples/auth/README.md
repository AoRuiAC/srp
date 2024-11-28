# SRP 鉴权样例

## 说明

SRP 使用 SSH 协议进行鉴权，支持自定义配置鉴权信息，这里给出一个从文件读入鉴权信息的样例。

该目录下已配置的鉴权情况如下：

提供了两个用户，配置如下：

- user 用户：
  - 允许以 `examples/auth/keys/user` 密钥登录到 Proxy 模块（在 `examples/auth/proxy_auth/user` 中定义）；
  - 允许对 `www.example.com`、`www.example-2.com` 两个网站进行代理（在 `examples/auth/proxy_rules/user` 中定义）；

- rpuser 用户：
  - 允许以 `examples/auth/keys/rpuser` 密钥登录到 ReverseProxy 模块（在 `examples/auth/reverseproxy_auth/rpuser` 中定义）；
  - 允许对 `www.example.com`、`www.example-3.com` 两个网站进行反向代理（在 `examples/auth/reverseproxy_rules/rpuser` 中定义）。

## 样例验证

### 服务启动

首先通过以下命令执行该样例的 SRP 服务：

```bash
go run ./examples/auth/main.go
```

### 反向代理连接

使用 user 身份进行反向代理的连接：

```bash
ssh -i examples/auth/keys/user -NR /www.example.com/80:127.0.0.1:8000 -p 8022 user@127.0.0.1
```

会得到连接建立失败的提示，服务端和客户端分别得到以下错误：

```
# server
INFO User user is not allowed to handle reverse proxy request.

# ssh
Warning: remote port forwarding failed for listen path /www.example.com/80
```

使用 rpuser 身份进行反向代理的连接：

```bash
ssh -i examples/auth/keys/rpuser -NR /www.example.com/80:127.0.0.1:8000 -p 8022 rpuser@127.0.0.1
```

连接会建立成功，客户端没有任何提示，在服务端会得到以下日志：

```
INFO Forward request in /tmp/srp2191930279/www.example.com_80.sock is ready
```

再次尝试用 rpuser 身份对不允许的地址进行反向代理的连接：

```bash
ssh -i examples/auth/keys/rpuser -NR /www.example-2.com/80:127.0.0.1:8000 -p 8022 rpuser@127.0.0.1
```

会再次得到连接建立失败的提示，服务端和客户端分别得到以下错误：

```
# server
ERRO User rpuser request to proxy /www.example-2.com/80, but it's not allowed.

# ssh
Warning: remote port forwarding failed for listen path /www.example-2.com/80
```

### 代理连接

首先启动一个本地样例服务并用 rpuser 身份启动反向代理，将 `www.example.com:80` 和 `www.example-3.com:80` 代理到本地样例服务上：

```bash
# Shell 1
go run ./examples/http/main.go

# Shell 2
ssh -i examples/auth/keys/rpuser -N -R /www.example.com/80:127.0.0.1:8000 -R /www.example-3.com/80:127.0.0.1:8000 -p 8022 rpuser@127.0.0.1
```

接着使用 user 身份启动代理服务，将本地的 `127.0.0.1:8008` 代理到 `www.example.com:80` 上：

```bash
ssh -i examples/auth/keys/user -NL 127.0.0.1:8008:www.example.com:80 -p 8022 user@127.0.0.1
```

接着访问 `127.0.0.1:8008` 服务：

```bash
curl http://127.0.0.1:8008
```

可以得到样例服务返回的结果：

```
127.0.0.1:8008 /
```

接着使用 rpuser 身份重新启动代理服务，进行同样的代理：

```bash
ssh -i examples/auth/keys/rpuser -NL 127.0.0.1:8008:www.example.com:80 -p 8022 rpuser@127.0.0.1
```

同样进行 `curl` 请求之后会得到如下错误：

```
# server
ERRO Cannot create proxy for xxx: unauthenticated for proxy

# curl
curl: (52) Empty reply from server
```

接着尝试通过 `www.example-3.com` 来进行代理连接：

```bash
ssh -i examples/auth/keys/user -NL 127.0.0.1:8008:www.example.com-3:80 -p 8022 user@127.0.0.1 
```

进行 `curl` 请求之后依然会得到错误：

```
# server
ERRO Cannot create proxy for xxx: access denied

# curl
curl: (52) Empty reply from server
```

至此鉴权模块样例验证完成。
