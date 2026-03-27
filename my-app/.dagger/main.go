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

	telemetry "github.com/dagger/otel-go"
	"golang.org/x/sync/errgroup"
)

type MyApp struct {
	Src *dagger.Directory
}

func New(
	ws *dagger.Workspace,
) *MyApp {
	return &MyApp{
		Src: ws.Directory("/my-app", dagger.WorkspaceDirectoryOpts{Gitignore: true}),
	}
}

// +check
func (m *MyApp) BuildCheck(ctx context.Context) (*dagger.Container, error) {
	plt, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return nil, err
	}

	return m.build(ctx, plt)
}

func (m *MyApp) build(
	ctx context.Context,
	platform dagger.Platform,
) (_ *dagger.Container, rerr error) {
	ctx, span := Tracer().Start(ctx, "build: "+string(platform))
	defer telemetry.EndWithCause(span, &rerr)

	build := dag.Container().
		From("node:latest").
		WithWorkdir("/app").
		WithDirectory(".", m.Src).
		WithExec([]string{"npm", "install"}).
		WithEnvVariable("PUBLIC_BACKEND_PORT", "30001").
		WithExec([]string{"npm", "run", "build"}).
		Directory("build")

	return dag.Container(dagger.ContainerOpts{
		Platform: platform,
	}).
		From("node:latest").
		WithWorkdir("/app").
		WithFiles(".", []*dagger.File{m.Src.File("package.json"), m.Src.File("package-lock.json")}).
		WithExec([]string{"npm", "install"}).
		WithDirectory("build", build).
		WithEntrypoint([]string{"node", "build"}).Sync(ctx)
}

// +cache="never"
func (m *MyApp) Publish(
	ctx context.Context,
	tag string,
	registryPassword *dagger.Secret,
) error {
	plats := []dagger.Platform{"linux/arm64", "linux/amd64"}

	ctrs := make([]*dagger.Container, len(plats))

	errg, gctx := errgroup.WithContext(ctx)

	for i, plt := range plats {
		errg.Go(func() error {
			build, err := m.build(gctx, plt)
			if err != nil {
				return err
			}

			ctrs[i] = build

			return nil
		})
	}

	if err := errg.Wait(); err != nil {
		return err
	}

	_, err := dag.Container().
		WithRegistryAuth("ghcr.io", "USERNAME", registryPassword).
		Publish(ctx, fmt.Sprintf("ghcr.io/chrisjpalmer/shoppinglist:frontend-%s", tag), dagger.ContainerPublishOpts{
			PlatformVariants: ctrs,
		})

	if err != nil {
		return err
	}

	return nil
}
