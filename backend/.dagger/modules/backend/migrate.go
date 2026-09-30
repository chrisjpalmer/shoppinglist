package main

import (
	"context"
)

const schemaPath = "backend/sql/schema.sql"

const atlasVersion = "1.1.6"

// MigrateCheck - checks whether the previous schema on the master branch
// can be successfully migrated to the new schema
// +check
func (m *Backend) MigrateCheck(ctx context.Context) error {
	prevsql := m.RootSrc.AsGit().Branch("master").Tree().File(schemaPath)

	newsql := m.RootSrc.File(schemaPath)

	return dag.MigrationTools().CheckMigrationValid(ctx, prevsql, newsql)
}
