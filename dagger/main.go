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
	container, err := m.BuildEnv(ctx, source, luaVersion, imageName, imageVersion, luarocksVersion)
	if err != nil {
		return "", err
	}

	lv, _ := NewLuaVersion(luaVersion)
	return container.WithExec([]string{lv.Executable(), "httpbin.lua"}).Stdout(ctx)
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

	iv := NewImageVersion(imageName, imageVersion)
	lr := NewLuarocksVersion(luarocksVersion)

	container := dag.Container().
		From(iv.ImageName()).
		WithMountedCache("/var/cache/apk", dag.CacheVolume("apk-cache")).
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "upgrade"}).
		WithExec([]string{
			"apk", "add", "--no-cache",
			"make", "tar", "wget", "gcc", "libc-dev", "openssl-dev",
			lv.PackageName(), lv.DevPackageName(),
		}).
		WithWorkdir("/src").
		WithExec([]string{"wget", lr.DownloadURL()}).
		WithExec([]string{"tar", "zxpf", lr.ArchiveName()}).
		WithWorkdir(lr.ExtractedDirPath()).
		WithExec([]string{"sh", "-c", lv.AssertSingleLuaH()}).
		WithWorkdir(lr.ExtractedDirPath()).
		WithExec([]string{"sh", "-c", strings.Join(lv.GetConfigureArgs(), " ")}).
		WithExec([]string{"make"}).
		WithExec([]string{"make"}).
		WithExec([]string{"make", "install"}).
		WithWorkdir("/").
		WithExec([]string{
			"rm", "-rf",
			lr.ExtractedDirPath(),
			lr.ArchiveName(),
		}).
		WithExec([]string{"luarocks", "install", "luasec"}).
		WithExec([]string{"luarocks", "install", "dkjson"})

	hollowbeakContainer := m.buildHollowbeak(source)

	return container.
		WithFile("/bin/hollowbeak", hollowbeakContainer.File("/src/hollowbeak/hollowbeak")).
		WithDirectory("/src", source).
		WithWorkdir("/src"), nil
}

func (m *Nobuffer) buildHollowbeak(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src/hollowbeak").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithEnvVariable("GOMODCACHE", "/go/pkg/mod").
		WithMountedCache("/go/build-cache", dag.CacheVolume("go-build-cache")).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "build", "-o", "hollowbeak"})
}
