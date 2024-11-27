package auth

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

type UserGlobs interface {
	Globs(ctx context.Context, user string) []glob.Glob
}

type UserGlobsFunc func(ctx context.Context, user string) []glob.Glob

func (f UserGlobsFunc) Globs(ctx context.Context, user string) []glob.Glob {
	return f(ctx, user)
}

type UserGlobsMap map[string][]glob.Glob

func (m UserGlobsMap) Globs(ctx context.Context, user string) []glob.Glob {
	return m[user]
}

type UserGlobsDir string

func (d UserGlobsDir) Globs(ctx context.Context, user string) []glob.Glob {
	filename := filepath.Join(string(d), user)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}
	sc := bufio.NewScanner(bytes.NewBuffer(data))
	ret := make([]glob.Glob, 0)
	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, ":") {
			line = line + ":*"
		}

		g, _ := glob.Compile(line, '.', ':', '/')
		if g != nil {
			ret = append(ret, g)
		}
	}
	return ret
}

func UserGlobsAuthorizer(c UserGlobs) Authorizer {
	return AuthorizeFunc(func(ctx context.Context, req AuthorizeRequest) bool {
		for _, g := range c.Globs(ctx, req.User) {
			if g.Match(req.Target) {
				return true
			}
		}
		return false
	})
}
