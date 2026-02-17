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
	"time"
)

type Backend struct {
	// +private
	Src *dagger.Directory
}

func New(
	// +defaultPath="/backend"
	src *dagger.Directory,
) *Backend {
	return &Backend{
		Src: src,
	}
}

// +cache="never"
func (m *Backend) PublishLinuxArm64(
	ctx context.Context,
	registryPassword *dagger.Secret,
) (string, error) {
	now := time.Now().Format("20060102-150405")

	build, err := m.BuildLinuxArm64(ctx)
	if err != nil {
		return "", err
	}

	_, err = build.
		WithRegistryAuth("ghcr.io", "USERNAME", registryPassword).
		Publish(ctx, fmt.Sprintf("ghcr.io/chrisjpalmer/shoppinglist:backend-%s", now))

	if err != nil {
		return "", err
	}

	return now, nil
}
