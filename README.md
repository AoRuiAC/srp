# SRP

SRP 是一个基于 SSH 协议开发的安全反向代理工具。通过 SRP，你可以将本地服务代理到其他的公开或私有网络中。

SRP 完全兼容 OpenSSH 客户端，你可以根据自己的需求和使用习惯自由选择使用 OpenSSH 客户端，或是 SRP 客户端进行连接和代理。

## 使用说明

### SRP 服务端

TODO

### OpenSSH 客户端

#### 基础的代理能力

可以通过 OpenSSH 连接 SRP 服务器进行反向代理，唯一需要注意的地方是，你必须按照一定格式填写代理的地址，示例如下：

```bash
ssh -NR /www.example.com/80:127.0.0.1:8000 SERVER_ADDR
```

通过以上命令，可以连接 SRP 服务器并进行代理，将 `www.example.com:80` 代理到本地的 `8000` 端口。

完成了反向代理之后，并不意味着在服务端可以通过 `www.example.com:80` 来访问代理的目标服务，需要在另一个本地客户端开启代理，示例如下：

```bash
ssh -NL 127.0.0.1:8000:www.example.com:80 SERVER_ADDR
```

通过以上命令，连接服务器并进行代理后，即可将本地的 `8000` 端口通过 `www.example.com:80` 代理到第一个客户端所在位置的 `8000` 端口。

#### 动态转发代理

从以上的使用来看，SRP 似乎并没有比 OpenSSH 服务器提供更丰富的代理能力？

其中的奥秘正在于中转时所使用的 `www.example.com:80`，允许将服务反向代理在形如 `www.example.com:80` 的主机名+端口上，而不是代理在 SSH 服务器的物理端口上，意味着代理服务的数量不再受限——当然这并不重要。最大的意义在于它提供了通过主机名+端口来访问所代理的服务的方式。通过 OpenSSH 的 DynamicForward 功能可以更好地体验到这种代理带来的便捷。

通过以下命令启动一个带有 DynamicForward 的 SSH 连接：

```bash
ssh -ND 127.0.0.1:1035 SERVER_ADDR
```

执行以上命令之后，OpenSSH 客户端会在本地的 `1035` 端口运行一个 socks5 服务，通过这个服务，可以用 `www.example.com:80` 这个地址来访问最终代理的服务。可以用 `curl` 命令通过以下方式进行验证：

```bash
curl --proxy socks5h://127.0.0.1:1035 http://www.example.com/
```

除此之外，你还可以借助 SwitchyOmega 等插件，在浏览器中直接通过 `http://www.example.com/` 访问到代理的目标服务。

#### SSH ProxyJump

当你对一个 SSH 服务进行反向代理时，你还可以将 SRP 服务当作一个 SSH 跳板机使用，并且使用起来非常便捷。

例如你将 SRP 服务的 `github.com:22` 代理到某个目标 SSH 服务上，但你要访问时，可以直接通过以下命令进行代理访问：

```bash
ssh -J SERVER_ADDR github.com
```

或在 SSH 配置文件中加入以下配置，可以达到同样的效果：

```
Host github.com
    ProxyJump SERVER_ADDR
```

### SRP 客户端

TODO