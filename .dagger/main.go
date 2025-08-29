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

	"gopkg.in/yaml.v3"
)

type Shoppinglist struct{}

func (m *Shoppinglist) Deploy(ctx context.Context,
	// +defaultPath="/"
	// +ignore=["/local"]
	src *dagger.Directory,
	registryPassword *dagger.Secret,
	kubectlFile *dagger.Secret,
	clusterIP string,
	clusterHost string,
) error {
	if err := m.DeployBackend(ctx, src, registryPassword, kubectlFile, clusterIP, clusterHost); err != nil {
		return err
	}

	if err := m.DeployFrontend(ctx, src, registryPassword, kubectlFile, clusterIP, clusterHost); err != nil {
		return err
	}

	return nil
}

func (m *Shoppinglist) DeployBackend(ctx context.Context,
	// +defaultPath="/"
	// +ignore=["/local"]
	src *dagger.Directory,
	registryPassword *dagger.Secret,
	kubectlFile *dagger.Secret,
	clusterIP string,
	clusterHost string,
) error {
	backend := src.Directory("backend")

	tag, err := dag.Backend().PublishLinuxArm64(
		ctx,
		registryPassword,
		dagger.BackendPublishLinuxArm64Opts{Src: backend},
	)

	if err != nil {
		return err
	}

	valuesYaml, err := makeValuesYaml(tag)
	if err != nil {
		return err
	}

	helm, err := helm(ctx, clusterIP, clusterHost)
	if err != nil {
		return err
	}

	_, err = helm.
		Chart(backend.Directory("helm")).
		Package().
		WithKubeconfigSecret(kubectlFile).
		Install("backend", dagger.HelmPackageInstallOpts{Namespace: "backend", Values: []*dagger.File{valuesYaml}, CreateNamespace: true}).
		Name(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (m *Shoppinglist) DeployFrontend(ctx context.Context,
	// +defaultPath="/"
	// +ignore=["/local"]
	src *dagger.Directory,
	registryPassword *dagger.Secret,
	kubectlFile *dagger.Secret,
	clusterIP string,
	clusterHost string,
) error {
	frontend := src.Directory("my-app")

	tag, err := dag.MyApp().PublishLinuxArm64(
		ctx,
		registryPassword,
		dagger.MyAppPublishLinuxArm64Opts{Src: frontend},
	)

	if err != nil {
		return err
	}

	valuesYaml, err := makeValuesYaml(tag)
	if err != nil {
		return err
	}

	helm, err := helm(ctx, clusterIP, clusterHost)
	if err != nil {
		return err
	}

	_, err = helm.
		Chart(frontend.Directory("helm")).
		Package().
		WithKubeconfigSecret(kubectlFile).
		Install("frontend", dagger.HelmPackageInstallOpts{Namespace: "frontend", Values: []*dagger.File{valuesYaml}, CreateNamespace: true}).
		Name(ctx)

	if err != nil {
		return err
	}

	return nil
}

func helm(ctx context.Context, destIP, destHost string) (*dagger.Helm, error) {
	helm := dag.Container().From("alpine/helm:latest")

	etcHosts, err := helm.File("/etc/hosts").Contents(ctx)
	if err != nil {
		return nil, err
	}

	etcHosts += fmt.Sprintf("\n%s\t%s", destIP, destHost)

	helm = helm.WithNewFile("/etc/hosts", etcHosts)

	return dag.Helm(dagger.HelmOpts{
		Container: helm,
	}), nil
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
