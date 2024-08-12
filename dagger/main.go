package main

import (
	"context"
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
	lv := NewLuaVersion(luaVersion)
	return m.BuildTestEnv(source, lv, imageName, imageVersion, luarocksVersion).
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
) *dagger.Container {
	lv := NewLuaVersion(luaVersion)
	return m.baseEnv(source, lv, imageName, imageVersion, luarocksVersion)
}

func (m *Nobuffer) BuildTestEnv(
	source *dagger.Directory,
	lv LuaVersion,
	// +optional
	imageName string,
	// +optional
	imageVersion string,
	// +optional
	luarocksVersion string,
) *dagger.Container {
	return m.installTestDependencies(
		m.baseEnv(source, lv, imageName, imageVersion, luarocksVersion),
		lv,
	)
}

func (m *Nobuffer) baseEnv(
	source *dagger.Directory,
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

	return dag.Container().
		From(iv.ImageName()).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "upgrade"}).
		WithExec([]string{
			"apk", "add", "--no-cache",
			"make",
			"tar",
			"wget",
			lv.PackageName(),
			lv.DevPackageName(),
		}).
		WithExec([]string{"wget", lr.DownloadURL()}).
		WithExec([]string{"tar", "zxpf", lr.ArchiveName()}).
		WithWorkdir(lr.ExtractedDirPath()).
		WithExec([]string{"sh", "-c", strings.Join(lv.GetConfigureArgs(), " ")}).
		WithExec([]string{"make"}).
		WithExec([]string{"make", "install"}).
		WithWorkdir("/src").
		WithExec([]string{
			"rm", "-rf",
			lr.ExtractedDirPath(),
			lr.ArchiveName(),
		})
}

func (m *Nobuffer) installTestDependencies(
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
		WithExec([]string{"luarocks", "install", "luasocket"}).
		WithExec([]string{"luarocks", "install", "luasec"})
}
