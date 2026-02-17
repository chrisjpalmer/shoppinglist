package main

import (
	"context"
	"dagger/backend/internal/dagger"
	"fmt"

	"golang.org/x/mod/modfile"
)

// +check
func (m *Backend) BuildLinuxArm64(ctx context.Context) (*dagger.Container, error) {
	ctr, err := m.buildCtr(ctx)
	if err != nil {
		return nil, err
	}

	backend := ctr.
		WithEnvVariable("GOOS", "linux").
		WithEnvVariable("GOARCH", "arm64").
		WithExec([]string{"go", "build", "-o", "backend", "."}).
		File("backend")

	return dag.Container(dagger.ContainerOpts{
		Platform: "linux/arm64",
	}).From("alpine:latest").
		WithWorkdir("/app").
		WithFile("backend", backend).
		WithEntrypoint([]string{"./backend"}), nil
}

func (m *Backend) buildCtr(ctx context.Context) (*dagger.Container, error) {
	ver, err := m.goVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting go version: %w", err)
	}

	return dag.Container().
		From("golang:"+ver).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-cache")).
		WithEnvVariable("CGO_ENABLED", "0").
		WithWorkdir("/src").
		WithDirectory("/src", m.Src), nil

}

func (m *Backend) goVersion(ctx context.Context) (string, error) {
	s, err := m.Src.File("go.mod").Contents(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting file contents: %w", err)
	}

	f, err := modfile.Parse("go.mod", []byte(s), nil)
	if err != nil {
		return "", fmt.Errorf("error parsing go.mod file: %w", err)
	}

	return f.Go.Version, nil
}
