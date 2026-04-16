package main

import (
	"context"
	"dagger/backend/internal/dagger"
	"fmt"
)

// GenerateProtos - generate protobuf codegen from .proto files
// +generate
func (m *Backend) GenerateProtos() *dagger.Changeset {
	gen := dag.Container().
		From("bufbuild/buf:1.65.0").
		WithExec([]string{"apk", "add", "--no-cache", "protobuf-dev"}).
		WithWorkdir("workdir").
		WithDirectory("proto", m.Src.Directory("proto")).
		WithFiles(".", []*dagger.File{m.Src.File("buf.yaml"), m.Src.File("buf.gen.yaml")}).
		WithExec([]string{"buf", "generate"}).
		Directory("gen")

	return m.Src.WithDirectory("gen", gen).Changes(m.Src)
}

// CheckProtos - check that the working tree's proto generated files are in sync.
// +check
func (m *Backend) CheckProtos(ctx context.Context) error {
	chgset := m.GenerateProtos()
	return assertEmpty(ctx, chgset)
}

// GenerateSqlc - generate sqlc codegen from .sql files
// +generate
func (m *Backend) GenerateSqlc() *dagger.Changeset {
	generated := dag.Container().
		From("sqlc/sqlc:1.29.0").
		WithWorkdir("workdir").
		WithDirectory("sql", m.Src.Directory("sql")).
		WithFile("sqlc.yaml", m.Src.File("sqlc.yaml")).
		WithExec([]string{"generate"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Directory("generated")

	return m.Src.WithDirectory("generated", generated).Changes(m.Src)
}

// CheckSqlc - check that the working tree's sqlc generated files are in sync.
// +check
func (m *Backend) CheckSqlc(ctx context.Context) error {
	chgset := m.GenerateSqlc()
	return assertEmpty(ctx, chgset)
}

// GenerateTempl - generate templ codegen from .templ files
// +generate
func (m *Backend) GenerateTempl(ctx context.Context) (*dagger.Changeset, error) {
	ctr, err := m.buildCtr(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting build container: %w", err)
	}

	gen := withTemplGenerate(ctr).Directory(".")

	return m.Src.WithDirectory(".", gen).Changes(m.Src), nil
}

// CheckTempl - check that the working tree's templ generated files are in sync.
// +check
func (m *Backend) CheckTempl(ctx context.Context) error {
	chgset, err := m.GenerateTempl(ctx)
	if err != nil {
		return fmt.Errorf("error generating templ sources: %w", err)
	}

	return assertEmpty(ctx, chgset)
}

func withTemplGenerate(ctr *dagger.Container) *dagger.Container {
	return ctr.WithExec([]string{"go", "tool", "templ", "generate"})
}

func assertEmpty(ctx context.Context, chgset *dagger.Changeset) error {
	empty, err := chgset.IsEmpty(ctx)
	if err != nil {
		return fmt.Errorf("error calling is empty: %w", err)
	}

	if !empty {
		return fmt.Errorf("generated files are out of sync")
	}

	return nil
}
