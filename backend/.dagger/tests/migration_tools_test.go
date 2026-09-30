package module_tests

import (
	"testing"

	dagger "github.com/shoppinglist/backend/.dagger/tests/internal/dagger/clients/migration-tools"
	"github.com/shoppinglist/backend/.dagger/tests/internal/dagger/clients/migration-tools/dag"

	_ "embed"
)

var (
	//go:embed golden/from.sql
	fromSql string
	//go:embed golden/from.sql
	toSql string
)

// TestMigrationToolsWithDB - tests that the migration tools work if the DB exists
func TestMigrationToolsWithDB(t *testing.T) {
	ctx := t.Context()

	_, err := toolsCtr(toSql).
		WithEnvVariable("DATABASE_FILE", "/app/local/local.db").
		WithFile("/app/local/local.db", dbForSchema(fromSql)).
		WithExec([]string{"/migrations/entrypoint.sh"}).
		Stdout(ctx)

	if err != nil {
		t.Fatal(err)
	}
}

// TestMigrationToolsNODB - tests that the migration tools work if the DB doesn't exist
func TestMigrationToolsNODB(t *testing.T) {
	ctx := t.Context()

	_, err := toolsCtr(toSql).
		WithEnvVariable("DATABASE_FILE", "/app/local/local.db").
		WithExec([]string{"/migrations/entrypoint.sh"}).
		Stdout(ctx)

	if err != nil {
		t.Fatal(err)
	}
}

// TestMigrationToolsNoDBEnv - tests that the migration tools correctly fail if the DATABASE_FILE var isn't present
func TestMigrationToolsNODBEnv(t *testing.T) {
	ctx := t.Context()

	code, err := toolsCtr(toSql).
		WithExec([]string{"/migrations/entrypoint.sh"}, dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny}).
		ExitCode(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if code != 1 {
		t.Fatal("expected entrypoint to return exit code 1 when DATABASE_FILE var was not set")
	}
}

func toolsCtr(toSql string) *dagger.Container {
	ctr := dag.Container().From("alpine:latest")

	ctr = dag.MigrationTools().InstallMigrationTools(ctr)

	return dag.MigrationTools().MountSchemaSQL(ctr, dag.File("to.sql", toSql))
}

// reads a schema file into a sqlite database and returns the database file produced
func dbForSchema(schema string) *dagger.File {
	return dag.Container().From("alpine/sqlite:3.51.2").
		WithFile("/app/schema.sql", dag.File("schema.sql", schema)).
		WithExec([]string{"/app/db.db", ".read /app/schema.sql"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		File("/app/db.db")
}
