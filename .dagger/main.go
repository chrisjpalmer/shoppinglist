// A generated module for Shoppinglist functions
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
	"dagger/shoppinglist/internal/dagger"
	"fmt"
	"time"

	telemetry "github.com/dagger/otel-go"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

type Shoppinglist struct {
	Backend  *dagger.Directory
	Frontend *dagger.Directory
}

func New(ws *dagger.Workspace) *Shoppinglist {
	src := ws.Directory("/", dagger.WorkspaceDirectoryOpts{Gitignore: true})

	return &Shoppinglist{
		Backend:  src.Directory("backend"),
		Frontend: src.Directory("my-app"),
	}
}

// +cache="never"
func (m *Shoppinglist) BuildAndDeploy(
	ctx context.Context,
	registryPassword *dagger.Secret,
	kubeEnv1 *dagger.Secret,
	kubeEnv2 *dagger.Secret,
) error {
	tag := time.Now().Format("20060102-150405")

	if err := m.Build(ctx, tag, registryPassword); err != nil {
		return fmt.Errorf("error while building: %w", err)
	}

	if err := m.Deploy(ctx, "env1", tag, kubeEnv1); err != nil {
		return fmt.Errorf("error while deploying to kubeEnv1: %w", err)
	}

	if err := m.Deploy(ctx, "env2", tag, kubeEnv2); err != nil {
		return fmt.Errorf("error while deploying to kubeEnv2: %w", err)
	}

	return nil
}

// +cache="never"
func (m *Shoppinglist) Deploy(
	ctx context.Context,
	env string,
	tag string,
	kubectlFile *dagger.Secret,
) (rerr error) {
	ctx, span := Tracer().Start(ctx, "deploy: "+env)
	defer telemetry.EndWithCause(span, &rerr)

	if err := m.deployBackend(ctx, tag, kubectlFile); err != nil {
		return err
	}

	if err := m.deployFrontend(ctx, tag, kubectlFile); err != nil {
		return err
	}

	return nil
}

// +cache="never"
func (m *Shoppinglist) Build(
	ctx context.Context,
	tag string,
	registryPassword *dagger.Secret,
) (rerr error) {
	ctx, span := Tracer().Start(ctx, "build")
	defer telemetry.EndWithCause(span, &rerr)

	errg, gctx := errgroup.WithContext(ctx)

	errg.Go(func() error { return m.publishBackend(gctx, tag, registryPassword) })

	errg.Go(func() error { return m.publishMigrateImage(gctx, tag, registryPassword) })

	errg.Go(func() error { return m.publishMyApp(gctx, tag, registryPassword) })

	if err := errg.Wait(); err != nil {
		return err
	}

	return nil
}

func (m *Shoppinglist) publishBackend(ctx context.Context, tag string, registryPassword *dagger.Secret) (rerr error) {
	ctx, span := Tracer().Start(ctx, "publish-backend")
	defer telemetry.EndWithCause(span, &rerr)

	return dag.Backend().Publish(ctx, tag, registryPassword)
}

func (m *Shoppinglist) publishMigrateImage(ctx context.Context, tag string, registryPassword *dagger.Secret) (rerr error) {
	ctx, span := Tracer().Start(ctx, "publish-migrate-image")
	defer telemetry.EndWithCause(span, &rerr)

	return dag.Backend().PublishMigrateImage(ctx, tag, registryPassword)
}

func (m *Shoppinglist) publishMyApp(ctx context.Context, tag string, registryPassword *dagger.Secret) (rerr error) {
	ctx, span := Tracer().Start(ctx, "publish-myapp")
	defer telemetry.EndWithCause(span, &rerr)

	return dag.MyApp().Publish(ctx, tag, registryPassword)
}

func (m *Shoppinglist) deployBackend(
	ctx context.Context,
	tag string,
	kubectlFile *dagger.Secret,
) (rerr error) {
	ctx, span := Tracer().Start(ctx, "deploy backend")
	defer telemetry.EndWithCause(span, &rerr)

	valuesYaml, err := makeValuesYaml(tag)
	if err != nil {
		return err
	}

	_, err = dag.Helm().
		Chart(m.Backend.Directory("helm")).
		Package().
		WithKubeconfigSecret(kubectlFile).
		Install("backend", dagger.HelmPackageInstallOpts{Namespace: "backend", Values: []*dagger.File{valuesYaml}, CreateNamespace: true}).
		Name(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (m *Shoppinglist) deployFrontend(
	ctx context.Context,
	tag string,
	kubectlFile *dagger.Secret,
) (rerr error) {
	ctx, span := Tracer().Start(ctx, "deploy frontend")
	defer telemetry.EndWithCause(span, &rerr)

	valuesYaml, err := makeValuesYaml(tag)
	if err != nil {
		return err
	}

	_, err = dag.Helm().
		Chart(m.Frontend.Directory("helm")).
		Package().
		WithKubeconfigSecret(kubectlFile).
		Install("frontend", dagger.HelmPackageInstallOpts{Namespace: "frontend", Values: []*dagger.File{valuesYaml}, CreateNamespace: true}).
		Name(ctx)

	if err != nil {
		return err
	}

	return nil
}

func makeValuesYaml(tag string) (*dagger.File, error) {
	type Image struct {
		Tag string `yaml:"tag"`
	}
	type ValuesYaml struct {
		Image Image `yaml:"image"`
	}

	out, err := yaml.Marshal(ValuesYaml{Image: Image{Tag: tag}})
	if err != nil {
		return nil, err
	}

	return dag.File("values.yaml", string(out)), nil
}
