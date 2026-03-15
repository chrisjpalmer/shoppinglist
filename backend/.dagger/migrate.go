package main

import (
	"context"
	"dagger/backend/internal/dagger"
)

const schemaPath = "backend/sql/schema.sql"

const atlasVersion = "arigaio/atlas:1.1.6-extended-alpine"

// MigrateLocal - migrates the passed in database and returns it
func (m *Backend) MigrateLocal(ctx context.Context, localdb *dagger.File) *dagger.File {
	newsql := m.RootSrc.File(schemaPath)

	return migrate(localdb, newsql).File("prev.db")
}

// MigrateCheck - checks whether the previous schema on the master branch
// can be successfully migrated to the new schema
// +check
func (m *Backend) MigrateCheck(ctx context.Context) error {
	prev := m.RootSrc.AsGit().Branch("master").Tree().File(schemaPath)
	prevdb := dbForSchema(prev)

	newsql := m.RootSrc.File(schemaPath)

	_, err := migrate(prevdb, newsql).Stdout(ctx)

	if err != nil {
		return err
	}

	return nil
}

func migrate(prevdb, newsql *dagger.File) *dagger.Container {
	return dag.Container().
		From(atlasVersion).
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

func dbForSchema(schema *dagger.File) *dagger.File {
	return dag.Container().From("alpine/sqlite:3.51.2").
		WithFile("/app/schema.sql", schema).
		WithExec([]string{"/app/db.db", ".read /app/schema.sql"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		File("/app/db.db")
}
