package main

import (
	"context"
	"dagger/backend/internal/dagger"
	"errors"
	"fmt"
)

// +cache="never"
func (m *Backend) PublishMigrateImage(
	ctx context.Context,
	tag string,
	registryPassword *dagger.Secret,
) error {
	plats := []dagger.Platform{"linux/arm64", "linux/amd64"}

	var ctrs []*dagger.Container

	for _, plt := range plats {
		ctrs = append(ctrs, m.migrateImage(plt))
	}

	_, err := dag.Container().
		WithRegistryAuth("ghcr.io", "USERNAME", registryPassword).
		Publish(ctx, fmt.Sprintf("ghcr.io/chrisjpalmer/shoppinglist:migrate-%s", tag), dagger.ContainerPublishOpts{
			PlatformVariants: ctrs,
		})

	if err != nil {
		return err
	}

	return nil
}

func (m *Backend) migrateImage(platform dagger.Platform) *dagger.Container {
	new := m.RootSrc.File(schemaPath)

	entrypoint := dag.CurrentModule().Source().File("migrate_image/entrypoint.sh")

	return dag.Container(dagger.ContainerOpts{Platform: platform}).
		From(atlasVersion).
		WithWorkdir("/app").
		WithFile("new.sql", new).
		WithFile("/entrypoint.sh", entrypoint, dagger.ContainerWithFileOpts{Permissions: 444}).
		WithEntrypoint([]string{"sh", "-c", "/entrypoint.sh"})
}

// TestMigrateImageWithDB - tests that the migrate image works if the DB exists
// +check
func (m *Backend) TestMigrateImageWithDB(ctx context.Context) error {
	schema := dag.CurrentModule().Source().File("golden/schema.sql")

	plt, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return err
	}

	_, err = m.migrateImage(plt).
		WithEnvVariable("DATABASE_FILE", "/app/local/local.db").
		WithFile("/app/local/local.db", dbForSchema(schema)).
		WithExec(nil, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)

	return err
}

// TestMigrateImageNODB - tests that the migrate image works if the DB exists
// +check
func (m *Backend) TestMigrateImageNODB(ctx context.Context) error {
	plt, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return err
	}

	_, err = m.migrateImage(plt).
		WithEnvVariable("DATABASE_FILE", "/app/local/local.db").
		WithExec(nil, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)

	return err
}

// TestMigrateImageNoDBEnv - tests that the migrate image correctly fails if the DATABASE_FILE var isn't present
// +check
func (m *Backend) TestMigrateImageNODBEnv(ctx context.Context) error {
	plt, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return err
	}

	code, err := m.migrateImage(plt).
		WithExec(nil, dagger.ContainerWithExecOpts{UseEntrypoint: true, Expect: dagger.ReturnTypeFailure}).
		ExitCode(ctx)

	if err != nil {
		return err
	}

	if code != 1 {
		return errors.New("expected entrypoint to return exit code 1 when DATABASE_FILE var was not set")
	}

	return nil
}
