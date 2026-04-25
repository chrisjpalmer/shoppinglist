// A generated module for Frontend functions
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
	"dagger/frontend/internal/dagger"
	"fmt"
	"strconv"
	"strings"

	telemetry "github.com/dagger/otel-go"
	"golang.org/x/sync/errgroup"
)

const (
	// helmBackendPort - the port of the api server from the backend helm chart.
	// This value must be synced with the helm chart
	helmBackendPort = 30001
)

type Frontend struct {
	// +private
	RootSrc *dagger.Directory
	Src     *dagger.Directory
}

func New(
	ws *dagger.Workspace,
) *Frontend {
	return &Frontend{
		RootSrc: ws.Directory("/", dagger.WorkspaceDirectoryOpts{Gitignore: true}),
		Src:     ws.Directory("/frontend", dagger.WorkspaceDirectoryOpts{Gitignore: true}),
	}
}

// +check
func (m *Frontend) BuildCheck(ctx context.Context) (*dagger.Container, error) {
	plt, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return nil, err
	}

	return m.build(ctx, plt, helmBackendPort)
}

func (m *Frontend) build(
	ctx context.Context,
	platform dagger.Platform,
	backendPort int,
) (_ *dagger.Container, rerr error) {
	ctx, span := Tracer().Start(ctx, "build: "+string(platform))
	defer telemetry.EndWithCause(span, &rerr)

	nodeVer, err := m.nodeVersion(ctx)
	if err != nil {
		return nil, err
	}

	build := m.buildCtr(nodeVer).
		WithExec([]string{"npm", "install"}).
		WithEnvVariable("PUBLIC_BACKEND_PORT", strconv.Itoa(backendPort)).
		WithExec([]string{"npm", "run", "build"}).
		Directory("build")

	return dag.Container(dagger.ContainerOpts{
		Platform: platform,
	}).
		From("node:" + nodeVer).
		WithMountedCache("/root/.npm", dag.CacheVolume("npm-build-cache")).
		WithWorkdir("/app").
		WithFiles(".", []*dagger.File{m.Src.File("package.json"), m.Src.File("package-lock.json")}).
		WithExec([]string{"npm", "install", "--omit=dev"}).
		WithDirectory("build", build).
		WithEntrypoint([]string{"node", "build"}).Sync(ctx)
}

func (m *Frontend) buildCtr(nodeVer string) *dagger.Container {
	return dag.Container().
		From("node:" + nodeVer).
		WithMountedCache("/root/.npm", dag.CacheVolume("npm-build-cache")).
		WithWorkdir("/app").
		WithDirectory(".", m.RootSrc).
		WithWorkdir("frontend")
}

func (m *Frontend) nodeVersion(ctx context.Context) (string, error) {
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

// +cache="never"
func (m *Frontend) Publish(
	ctx context.Context,
	tag string,
	registryPassword *dagger.Secret,
) error {
	plats := []dagger.Platform{"linux/arm64", "linux/amd64"}

	ctrs := make([]*dagger.Container, len(plats))

	errg, gctx := errgroup.WithContext(ctx)

	for i, plt := range plats {
		errg.Go(func() error {
			build, err := m.build(gctx, plt, helmBackendPort)
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
