package main

import (
	"context"
	"dagger/backend/internal/dagger"
	"errors"
)

// withMigrationTools - installs the migration tools into the specified container
func (m *Backend) withMigrationTools(ctr *dagger.Container) *dagger.Container {
	entrypoint := dag.CurrentModule().Source().File("migration_tools/entrypoint.sh")

	return ctr.
		WithExec([]string{"apk", "add", "curl"}).
		WithEnvVariable("ATLAS_VERSION", "v"+atlasVersion).
		WithExec([]string{"sh", "-c", "curl -sSf https://atlasgo.sh | sh"}).
		WithExec([]string{"apk", "del", "curl"}).
		WithFile("/migrations/entrypoint.sh", entrypoint, dagger.ContainerWithFileOpts{Permissions: 444})
}

// withMigrationSQL - installs the migration sql into the expected location
func (m *Backend) withMigrationSQL(ctr *dagger.Container, toSql *dagger.File) *dagger.Container {
	return ctr.WithFile("/migrations/to.sql", toSql)
}

// TestMigrationToolsWithDB - tests that the migration tools work if the DB exists
// +check
func (m *Backend) TestMigrationToolsWithDB(ctx context.Context) error {
	curMod := dag.CurrentModule().Source()

	fromSql := curMod.File("golden/from.sql")

	toSql := curMod.File("golden/to.sql")

	_, err := m.testMigrationTools(toSql).
		WithEnvVariable("DATABASE_FILE", "/app/local/local.db").
		WithFile("/app/local/local.db", dbForSchema(fromSql)).
		WithExec([]string{"/migrations/entrypoint.sh"}).
		Stdout(ctx)

	return err
}

// TestMigrationToolsNODB - tests that the migration tools work if the DB doesn't exist
// +check
func (m *Backend) TestMigrationToolsNODB(ctx context.Context) error {
	toSql := m.RootSrc.File(schemaPath)

	_, err := m.testMigrationTools(toSql).
		WithEnvVariable("DATABASE_FILE", "/app/local/local.db").
		WithExec([]string{"/migrations/entrypoint.sh"}).
		Stdout(ctx)

	return err
}

// TestMigrationToolsNoDBEnv - tests that the migration tools correctly fail if the DATABASE_FILE var isn't present
// +check
func (m *Backend) TestMigrationToolsNODBEnv(ctx context.Context) error {
	toSql := m.RootSrc.File(schemaPath)

	code, err := m.testMigrationTools(toSql).
		WithExec([]string{"/migrations/entrypoint.sh"}, dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny}).
		ExitCode(ctx)

	if err != nil {
		return err
	}

	if code != 1 {
		return errors.New("expected entrypoint to return exit code 1 when DATABASE_FILE var was not set")
	}

	return nil
}

func (m *Backend) testMigrationTools(toSql *dagger.File) *dagger.Container {
	ctr := dag.Container().From("alpine:latest")

	ctr = m.withMigrationTools(ctr)

	return m.withMigrationSQL(ctr, toSql)
}
