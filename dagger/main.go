package main

import (
	"context"

	"dagger/nobuffer/internal/dagger"
)

type Nobuffer struct{}

// Return the result of running unit tests
func (m *Nobuffer) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.BuildEnv(source).
		WithExec([]string{"lua5.4", "httpbin.lua"}).
		Stdout(ctx)
}

// Build a ready-to-use development environment
func (m *Nobuffer) BuildEnv(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"/bin/sh", "luarocks-install.sh"}).
		WithExec([]string{"luarocks", "install", "luasec"}).
		WithExec([]string{"luarocks", "install", "luasocket"})
}
