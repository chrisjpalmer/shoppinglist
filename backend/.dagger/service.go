package main

import (
	"context"
	"dagger/backend/internal/dagger"
	"fmt"
)

const (
	// localBackendApiPort - the default port that the backend api server
	// is served from.
	localBackendApiPort = 8080

	// localBackendShoppingSitePort - the default port that the shopping site
	// is served from.
	localBackendShoppingSitePort = 8081
)

// BackendService - runs the backend service inside a container
// +up
func (m *Backend) BackendService(ctx context.Context, ws *dagger.Workspace) (*dagger.Service, error) {
	plt, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return nil, err
	}

	ctr, err := m.build(ctx, plt)
	if err != nil {
		return nil, fmt.Errorf("error building backend: %w", err)
	}

	return ctr.
		WithMountedDirectory("local", ws.Directory("/backend/local")).
		WithExposedPort(localBackendApiPort).
		WithExposedPort(localBackendShoppingSitePort).
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true}), nil
}
