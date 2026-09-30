package main

import (
	"context"
	"dagger/migration-tools/internal/dagger"
	"fmt"
)

const atlasVersion = "1.1.6"

type MigrationTools struct {
}

// MigrateDatabase - migrates the passed in database, to the provided schema and returns it
func (m *MigrationTools) MigrateDatabase(ctx context.Context, db *dagger.File, schemaSQL *dagger.File) *dagger.File {
	return migrate(db, schemaSQL).File("prev.db")
}

// CheckMigrationValid - checks whether a migration from the previous schema to the new schema succeeds.
func (m *MigrationTools) CheckMigrationValid(ctx context.Context, prevSchemaSQL *dagger.File, newSchemaSQL *dagger.File) error {
	prevdb := dbForSchema(prevSchemaSQL)

	_, err := migrate(prevdb, newSchemaSQL).Stdout(ctx)

	if err != nil {
		return fmt.Errorf("failed to migrate from previous to new schema: %w", err)
	}

	return nil
}

func migrate(prevdb, newsql *dagger.File) *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("arigaio/atlas:%s-extended-alpine", atlasVersion)).
		WithWorkdir("/app").
		WithFile("prev.db", prevdb).
		WithFile("new.sql", newsql).
		WithExec([]string{
			"schema", "apply",
			"--url", "sqlite:///app/prev.db",
			"--to", "file:///app/new.sql",
			"--dev-url", "sqlite://dev?mode=memory",
			"--auto-approve",
		}, dagger.ContainerWithExecOpts{UseEntrypoint: true})
}

// InstallMigrationTools - installs the migration tools into the specified container
func (m *MigrationTools) InstallMigrationTools(ctr *dagger.Container) *dagger.Container {
	entrypoint := dag.CurrentModule().Source().File("entrypoint.sh")

	return ctr.
		WithExec([]string{"apk", "add", "curl"}).
		WithEnvVariable("ATLAS_VERSION", "v"+atlasVersion).
		WithExec([]string{"sh", "-c", "curl -sSf https://atlasgo.sh | sh"}).
		WithExec([]string{"apk", "del", "curl"}).
		WithFile("/migrations/entrypoint.sh", entrypoint, dagger.ContainerWithFileOpts{Permissions: 444})
}

// MountSchemaSQL - installs the migration sql into the expected location
func (m *MigrationTools) MountSchemaSQL(ctr *dagger.Container, schemaSQL *dagger.File) *dagger.Container {
	return ctr.WithFile("/migrations/to.sql", schemaSQL)
}

// reads a schema file into a sqlite database and returns the database file produced
func dbForSchema(schema *dagger.File) *dagger.File {
	return dag.Container().From("alpine/sqlite:3.51.2").
		WithFile("/app/schema.sql", schema).
		WithExec([]string{"/app/db.db", ".read /app/schema.sql"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		File("/app/db.db")
}
