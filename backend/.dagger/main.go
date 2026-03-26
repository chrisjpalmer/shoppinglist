// A generated module for Backend functions
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
	"dagger/backend/internal/dagger"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Backend struct {
	// +private
	RootSrc *dagger.Directory
	// +private
	Src *dagger.Directory
}

func New(
	ws *dagger.Workspace,
) *Backend {
	return &Backend{
		RootSrc: ws.Directory("/", dagger.WorkspaceDirectoryOpts{Gitignore: true}),
		Src:     ws.Directory("/backend", dagger.WorkspaceDirectoryOpts{Gitignore: true}),
	}
}

// +cache="never"
func (m *Backend) Publish(
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
		Publish(ctx, fmt.Sprintf("ghcr.io/chrisjpalmer/shoppinglist:backend-%s", tag), dagger.ContainerPublishOpts{
			PlatformVariants: ctrs,
		})

	if err != nil {
		return err
	}

	return nil
}
