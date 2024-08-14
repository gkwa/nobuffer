package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/nobuffer/internal/dagger"
)

type Nobuffer struct{}

func (m *Nobuffer) Test(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	luaVersion string,
	// +optional
	imageName string,
	// +optional
	imageVersion string,
	// +optional
	luarocksVersion string,
) (string, error) {
	lv, err := NewLuaVersion(luaVersion)
	if err != nil {
		return "", fmt.Errorf("failed to create LuaVersion: %w", err)
	}
	return m.BuildTestEnv(ctx, source, lv.version, imageName, imageVersion, luarocksVersion).
		WithExec([]string{lv.Executable(), "httpbin.lua"}).
		Stdout(ctx)
}

func (m *Nobuffer) BuildEnv(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	luaVersion string,
	// +optional
	imageName string,
	// +optional
	imageVersion string,
	// +optional
	luarocksVersion string,
) (*dagger.Container, error) {
	lv, err := NewLuaVersion(luaVersion)

	if err != nil {
		return nil, fmt.Errorf("failed to create LuaVersion: %w", err)
	}
	baseContainer := m.baseEnv(lv, imageName, imageVersion, luarocksVersion)
	containerWithDeps := m.installDependencies(baseContainer, lv)
	return m.buildAndInstallHollowbeak(containerWithDeps, source), nil
}

func (m *Nobuffer) BuildTestEnv(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	luaVersion string,
	// +optional
	imageName string,
	// +optional
	imageVersion string,
	// +optional
	luarocksVersion string,
) *dagger.Container {
	lv, err := NewLuaVersion(luaVersion)
	if err != nil {
		panic(err)
	}
	baseContainer := m.baseEnv(lv, imageName, imageVersion, luarocksVersion)
	containerWithDeps := m.installDependencies(baseContainer, lv)
	return m.buildAndInstallHollowbeak(containerWithDeps, source)
}

func (m *Nobuffer) baseEnv(
	lv LuaVersion,
	// +optional
	imageName string,
	// +optional
	imageVersion string,
	// +optional
	luarocksVersion string,
) *dagger.Container {
	iv := NewImageVersion(imageName, imageVersion)
	lr := NewLuarocksVersion(luarocksVersion)



	base := dag.Container().
		From(iv.ImageName()).
		WithMountedCache("/var/cache/apk", dag.CacheVolume("apk-cache")).
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "upgrade"}).
		WithExec([]string{
			"apk", "add", "--no-cache",
			"make",
			"tar",
			"wget",
			lv.PackageName(),
			lv.DevPackageName(),
		})

	luaCacheVolume := dag.CacheVolume(fmt.Sprintf("lua-cache-%s", lv.version))
	luarocksCacheVolume := dag.CacheVolume(fmt.Sprintf("luarocks-cache-%s", lr.version))

	cont := base.
		WithMountedCache("/usr/local/lib/lua", luaCacheVolume).
		WithMountedCache("/usr/local/share/lua", luaCacheVolume).
		WithMountedCache("/usr/local/lib/luarocks", luarocksCacheVolume).
		WithWorkdir("/").
		WithExec([]string{"wget", lr.DownloadURL()}).
		WithExec([]string{"tar", "zxpf", lr.ArchiveName()}).
		WithWorkdir(fmt.Sprintf("/luarocks-%s", lr.version)).
		WithExec([]string{"sh", "-c", lv.AssertSingleLuaH()}).
		WithExec([]string{"sh", "-c", strings.Join(lv.GetConfigureArgs(), " ")}).
		WithExec([]string{"make"}).
		WithExec([]string{"make", "install"}).
		WithWorkdir("/").
		WithExec([]string{
			"rm", "-rf",
			fmt.Sprintf("/luarocks-%s", lr.version),
			lr.ArchiveName(),
		})

	return cont
}

func (m *Nobuffer) installDependencies(
	base *dagger.Container,
	lv LuaVersion,
) *dagger.Container {
	return base.
		WithExec([]string{
			"apk", "add", "--no-cache",
			"gcc",
			"libc-dev",
			"openssl-dev",
		}).
		WithExec([]string{"luarocks", "install", "luasec"}).
		WithExec([]string{"luarocks", "install", "dkjson"})
}

func (m *Nobuffer) buildAndInstallHollowbeak(
	base *dagger.Container,
	source *dagger.Directory,
) *dagger.Container {
	builder := dag.Container().
		From("golang:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src/hollowbeak").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "build", "-o", "hollowbeak"})

	return base.
		WithFile("/bin/hollowbeak", builder.File("/src/hollowbeak/hollowbeak")).
		WithDirectory("/src", source).
		WithWorkdir("/src")
}
