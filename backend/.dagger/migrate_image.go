package main

import (
	"context"
	"dagger/backend/internal/dagger"
	"errors"
	"fmt"
)

// PublishMigrateImage - builds an image that contains the atlas migration tool
// as well as the new schema that the database needs to be migrated to
// +cache="never"
func (m *Backend) PublishMigrateImage(
	ctx context.Context,
	tag string,
	registryPassword *dagger.Secret,
) error {
	plats := []dagger.Platform{"linux/arm64", "linux/amd64"}

	toSql := m.RootSrc.File(schemaPath)

	var ctrs []*dagger.Container

	for _, plt := range plats {
		ctrs = append(ctrs, m.migrateImage(plt, toSql))
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

func (m *Backend) migrateImage(platform dagger.Platform, toSql *dagger.File) *dagger.Container {

	entrypoint := dag.CurrentModule().Source().File("migrate_image/entrypoint.sh")

	return dag.Container(dagger.ContainerOpts{Platform: platform}).
		From(atlasVersion).
		WithWorkdir("/app").
		WithFile("to.sql", toSql).
		WithFile("/entrypoint.sh", entrypoint, dagger.ContainerWithFileOpts{Permissions: 444}).
		WithEntrypoint([]string{"sh", "-c", "/entrypoint.sh"})
}

// TestMigrateImageWithDB - tests that the migrate image works if the DB exists
// +check
func (m *Backend) TestMigrateImageWithDB(ctx context.Context) error {
	curMod := dag.CurrentModule().Source()

	fromSql := curMod.File("golden/from.sql")

	toSql := curMod.File("golden/to.sql")

	plt, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return err
	}

	_, err = m.migrateImage(plt, toSql).
		WithEnvVariable("DATABASE_FILE", "/app/local/local.db").
		WithFile("/app/local/local.db", dbForSchema(fromSql)).
		WithExec(nil, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)

	return err
}

// TestMigrateImageNODB - tests that the migrate image works if the DB exists
// +check
func (m *Backend) TestMigrateImageNODB(ctx context.Context) error {
	toSql := m.RootSrc.File(schemaPath)

	plt, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return err
	}
	_, err = m.migrateImage(plt, toSql).
		WithEnvVariable("DATABASE_FILE", "/app/local/local.db").
		WithExec(nil, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Stdout(ctx)

	return err
}

// TestMigrateImageNoDBEnv - tests that the migrate image correctly fails if the DATABASE_FILE var isn't present
// +check
func (m *Backend) TestMigrateImageNODBEnv(ctx context.Context) error {
	toSql := m.RootSrc.File(schemaPath)

	plt, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return err
	}

	code, err := m.migrateImage(plt, toSql).
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
