package client

import gossh "golang.org/x/crypto/ssh"

type ProxyType int

const (
	DynamicForward ProxyType = iota
	LocalForward
	RemoteForward
)

type ProxyConfig struct {
	Type       ProxyType
	Network    string
	LocalHost  string
	LocalPort  string
	RemoteHost string
	RemotePort string
}

type ConnConfig struct {
	Network     string
	Address     string
	User        string
	AuthMethods []gossh.AuthMethod
	Proxies     []ProxyConfig
}
