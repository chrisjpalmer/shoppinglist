// A generated module for MyApp functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/my-app/internal/dagger"
	"fmt"
	"time"
)

type MyApp struct{}

func (m *MyApp) BuildLinuxArm64(
	// +defaultPath="/my-app"
	src *dagger.Directory,
) *dagger.Container {
	build := dag.Container().
		From("node:latest").
		WithWorkdir("/app").
		WithDirectory(".", src).
		WithExec([]string{"npm", "install"}).
		WithEnvVariable("PUBLIC_BACKEND_PORT", "30001").
		WithExec([]string{"npm", "run", "build"}).
		Directory("build")

	return dag.Container(dagger.ContainerOpts{
		Platform: "linux/arm64",
	}).
		From("node:latest").
		WithWorkdir("/app").
		WithFiles(".", []*dagger.File{src.File("package.json"), src.File("package-lock.json")}).
		WithExec([]string{"npm", "install"}).
		WithDirectory("build", build).
		WithEntrypoint([]string{"node", "build"})
}

func (m *MyApp) PublishLinuxArm64(
	ctx context.Context,
	// +defaultPath="/my-app"
	src *dagger.Directory,
	registryPassword *dagger.Secret,
) (string, error) {
	now := time.Now().Format("20060102-150405")

	_, err := m.BuildLinuxArm64(src).
		WithRegistryAuth("ghcr.io", "USERNAME", registryPassword).
		Publish(ctx, fmt.Sprintf("ghcr.io/chrisjpalmer/shoppinglist:frontend-%s", now))

	if err != nil {
		return "", err
	}

	return now, nil
}
