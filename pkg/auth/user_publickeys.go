package auth

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/ssh"
	gossh "golang.org/x/crypto/ssh"
)

type UserPublicKeys interface {
	PublicKeys(ctx context.Context, user string) []gossh.PublicKey
}

type UserPublicKeysFunc func(ctx context.Context, user string) []gossh.PublicKey

func (f UserPublicKeysFunc) PublicKeys(ctx context.Context, user string) []gossh.PublicKey {
	return f(ctx, user)
}

type UserPublicKeysMap map[string][]gossh.PublicKey

func (m UserPublicKeysMap) PublicKeys(ctx context.Context, user string) []gossh.PublicKey {
	return m[user]
}

type PublicKeysDir string

func (d PublicKeysDir) PublicKeys(ctx context.Context, user string) []gossh.PublicKey {
	filename := filepath.Join(string(d), user)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}
	sc := bufio.NewScanner(bytes.NewBuffer(data))
	ret := make([]gossh.PublicKey, 0)
	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		publickey, _, _, _, _ := gossh.ParseAuthorizedKey([]byte(line))
		if publickey != nil {
			ret = append(ret, publickey)
		}
	}
	return ret
}

func UserPublicKeysAuthenticator(c UserPublicKeys) Authenticator {
	return AuthenticateFunc(func(ctx context.Context, req AuthenticateRequest) bool {
		for _, publickey := range c.PublicKeys(ctx, req.User) {
			if ssh.KeysEqual(publickey, req.PublicKey) {
				return true
			}
		}
		return false
	})
}
