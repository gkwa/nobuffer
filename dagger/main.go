package main

import (
	"context"

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
	return m.BuildEnv(source, lv, imageName, imageVersion, luarocksVersion).
		WithExec([]string{lv.Executable(), "httpbin.lua"}).
		Stdout(ctx)
}

func (m *Nobuffer) BuildEnv(
	source *dagger.Directory,
	lv LuaVersion,
	// +optional
	imageName string,
	// +optional
	imageVersion string,
	// +optional
	luarocksVersion string,
) *dagger.Container {
	return m.InstallLuaDependencies(
		m.InstallLuarocks(
			m.InstallLua(
				m.BaseEnv(source, imageName, imageVersion),
				lv,
			),
			lv,
			luarocksVersion,
		),
		lv,
	)
}

func (m *Nobuffer) BaseEnv(
	source *dagger.Directory,
	// +optional
	imageName string,
	// +optional
	imageVersion string,
) *dagger.Container {
	iv := NewImageVersion(imageName, imageVersion)
	return dag.Container().
		From(iv.ImageName()).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "upgrade"}).
		WithExec([]string{
			"apk", "add", "--no-cache",
			"gcc",
			"libc-dev",
			"make",
			"openssl-dev",
			"readline-dev",
			"tar",
			"wget",
		})
}

func (m *Nobuffer) InstallLua(
	base *dagger.Container,
	lv LuaVersion,
) *dagger.Container {
	return base.
		WithExec([]string{
			"apk", "add", "--no-cache",
			lv.PackageName(),
			lv.DevPackageName(),
		})
}

func (m *Nobuffer) InstallLuarocks(
	base *dagger.Container,
	lv LuaVersion,
	// +optional
	luarocksVersion string,
) *dagger.Container {
	lr := NewLuarocksVersion(luarocksVersion)
	return base.
		WithExec([]string{
			"wget", lr.DownloadURL(),
		}).
		WithExec([]string{
			"tar", "zxpf", lr.ArchiveName(),
		}).
		WithWorkdir(lr.ExtractedDirPath()).
		WithExec([]string{
			"./configure",
			"--prefix=/usr",
			lv.LuaIncludePath(),
			"--with-lua=/usr",
			lv.InterpreterFlag(),
		}).
		WithExec([]string{"make"}).
		WithExec([]string{"make", "install"}).
		WithWorkdir("/src").
		WithExec([]string{
			"rm", "-rf",
			lr.ExtractedDirPath(),
			lr.ArchiveName(),
		})
}

func (m *Nobuffer) InstallLuaDependencies(
	base *dagger.Container,
	lv LuaVersion,
) *dagger.Container {
	return base.
		WithExec([]string{"luarocks", "install", "luasocket"}).
		WithExec([]string{"luarocks", "install", "luasec"})
}
