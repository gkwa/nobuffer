package main

import (
	"context"
	"fmt"

	"dagger/nobuffer/internal/dagger"
)

type Nobuffer struct{}

const (
	luaVersion      = "5.4"
	alpineVersion   = "3.20"
	luarocksVersion = "3.11.1"
)

// Return the result of running unit tests
func (m *Nobuffer) Test(
	ctx context.Context,
	source *dagger.Directory,
) (string, error) {
	return m.BuildEnv(source).
		WithExec([]string{fmt.Sprintf("lua%s", luaVersion), "httpbin.lua"}).
		Stdout(ctx)
}

// Build a ready-to-use development environment
func (m *Nobuffer) BuildEnv(
	source *dagger.Directory,
) *dagger.Container {
	return m.InstallLuaDependencies(m.InstallLuarocks(m.InstallLua(m.BaseEnv(source))))
}

func (m *Nobuffer) BaseEnv(
	source *dagger.Directory,
) *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("alpine:%s", alpineVersion)).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "upgrade"}).
		WithExec([]string{
			"apk", "add", "--no-cache",
			"wget",
			"tar",
			"gcc",
			"libc-dev",
			"make",
			"openssl-dev",
			"readline-dev",
		})
}

func (m *Nobuffer) InstallLua(
	base *dagger.Container,
) *dagger.Container {
	return base.
		WithExec([]string{
			"apk", "add", "--no-cache",
			fmt.Sprintf("lua%s", luaVersion),
			fmt.Sprintf("lua%s-dev", luaVersion),
		})
}

func (m *Nobuffer) InstallLuarocks(
	base *dagger.Container,
) *dagger.Container {
	return base.
		WithExec([]string{
			"wget", fmt.Sprintf("https://luarocks.org/releases/luarocks-%s.tar.gz", luarocksVersion),
		}).
		WithExec([]string{
			"tar", "zxpf", fmt.Sprintf("luarocks-%s.tar.gz", luarocksVersion),
		}).
		WithWorkdir(fmt.Sprintf("/src/luarocks-%s", luarocksVersion)).
		WithExec([]string{
			"./configure",
			"--prefix=/usr",
			fmt.Sprintf("--with-lua-include=/usr/include/lua%s", luaVersion),
			"--with-lua=/usr",
			fmt.Sprintf("--with-lua-interpreter=lua%s", luaVersion),
		}).
		WithExec([]string{"make"}).
		WithExec([]string{"make", "install"}).
		WithWorkdir("/src").
		WithExec([]string{
			"rm", "-rf",
			fmt.Sprintf("luarocks-%s", luarocksVersion),
			fmt.Sprintf("luarocks-%s.tar.gz", luarocksVersion),
		})
}

func (m *Nobuffer) InstallLuaDependencies(
	base *dagger.Container,
) *dagger.Container {
	return base.
		WithExec([]string{"luarocks", "install", "luasocket"}).
		WithExec([]string{"luarocks", "install", "luasec"})
}
