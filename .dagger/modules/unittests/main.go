package main

import (
	"context"
	"dagger/unittests/internal/dagger"
	"fmt"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// +collection
type Unittests struct {
	// +keys
	Paths []string
}

func New(
	ctx context.Context,
	ws *dagger.Workspace,
) (*Unittests, error) {
	paths, err := ws.Glob(ctx, "**/.dagger/tests")
	if err != nil {
		return nil, err
	}
	return &Unittests{
		Paths: paths,
	}, nil
}

// +get
func (m *Unittests) Package(path string) *Package {
	return &Package{Path: path}
}

type Package struct {
	Path string
}

// +check
func (p *Package) Test(ctx context.Context, ws *dagger.Workspace) error {
	ctr, err := p.testCtr(ctx, ws)
	if err != nil {
		return fmt.Errorf("error getting test container: %w")
	}

	_, err = ctr.WithExec([]string{"go", "test", "./..."}).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error running test: %w", err)
	}

	return nil
}

func (p *Package) testCtr(ctx context.Context, ws *dagger.Workspace) (*dagger.Container, error) {
	ver, err := p.goVersion(ctx, ws)
	if err != nil {
		return nil, fmt.Errorf("error getting go version: %w", err)
	}

	return dag.Container().
		From("golang:"+ver).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-cache")).
		WithEnvVariable("CGO_ENABLED", "0").
		WithDirectory("/src", ws.Directory("/")).
		WithWorkdir(filepath.Join("/src", p.Path)), nil
}

func (p *Package) goVersion(ctx context.Context, ws *dagger.Workspace) (string, error) {
	s, err := ws.File(filepath.Join(p.Path, "go.mod")).Contents(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting file contents: %w", err)
	}

	f, err := modfile.Parse("go.mod", []byte(s), nil)
	if err != nil {
		return "", fmt.Errorf("error parsing go.mod file: %w", err)
	}

	return f.Go.Version, nil
}
