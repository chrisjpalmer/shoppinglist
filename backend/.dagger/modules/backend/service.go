package main

import (
	"context"
	"dagger/backend/internal/dagger"
	"fmt"
)

const (
	// localBackendPort - the default port that the backend server
	// is served from.
	localBackendPort = 8080
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
		WithExposedPort(localBackendPort).
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true}), nil
}
