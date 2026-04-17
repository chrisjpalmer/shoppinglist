package main

import (
	"context"
	"dagger/frontend/internal/dagger"
	"fmt"
)

const (
	// localBackendPort - the port of the api server when the backend is served locally.
	// This value must be synced with the default port in the backend code
	localBackendPort = 8080

	// localFrontendPort - the port of the frontend website.
	// This is the default for the svelte setup being used.
	localFrontendPort = 3000
)

// FrontendService - runs the frontend service inside a container
// +up
func (m *Frontend) FrontendService(ctx context.Context) (*dagger.Service, error) {
	plt, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return nil, err
	}

	ctr, err := m.build(ctx, plt, localBackendPort)
	if err != nil {
		return nil, fmt.Errorf("error building frontend: %w", err)
	}

	return ctr.
		WithExposedPort(localFrontendPort).
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true}), nil
}
