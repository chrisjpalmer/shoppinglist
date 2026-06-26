package main

import (
	"context"
	"dagger/backend/internal/dagger"
	"fmt"
	"strings"
)

// GenerateTailwind generates the maintw.css file using tailwindcss
// +generate
func (m *Backend) GenerateTailwind(ctx context.Context) (*dagger.Changeset, error) {
	tw, err := m.tailwind(ctx)
	if err != nil {
		return nil, fmt.Errorf("error building tailwind container: %w", err)
	}

	maintw := tw.
		WithExec([]string{"npm", "run", "css"}).
		File("/src/shopping/assets/maintw.css")

	return m.Src.WithFile("shopping/assets/maintw.css", maintw).Changes(m.Src), nil
}

func (m *Backend) tailwind(ctx context.Context) (*dagger.Container, error) {
	nodeVer, err := m.nodeVersion(ctx)
	if err != nil {
		return nil, err
	}

	return dag.Container().
		From("node:"+nodeVer).
		WithMountedCache("/root/.npm", dag.CacheVolume("npm-tailwind-cache")).
		WithDirectory("/src", m.Src).
		WithWorkdir("/src/shopping/tailwind").
		WithExec([]string{"npm", "install"}), nil
}

func (m *Backend) nodeVersion(ctx context.Context) (string, error) {
	contents, err := m.RootSrc.File(".tool-versions").Contents(ctx)
	if err != nil {
		return "", fmt.Errorf("error reading .tool-versions: %w", err)
	}

	for line := range strings.SplitSeq(contents, "\n") {
		parts := strings.Fields(line)
		if len(parts) == 2 && parts[0] == "nodejs" {
			return parts[1], nil
		}
	}

	return "", fmt.Errorf("nodejs version not found in .tool-versions")
}
