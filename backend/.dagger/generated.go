package main

import (
	"context"
	"dagger/backend/internal/dagger"
	"fmt"
)

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

// Returns lines that match a pattern in the files of the provided Directory
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

// +check
func (m *Backend) CheckTempl(ctx context.Context) error {
	chgset, err := m.GenerateTempl(ctx)
	if err != nil {
		return fmt.Errorf("error generating templates: %w", err)
	}

	return assertEmpty(ctx, chgset)
}

// +generate
func (m *Backend) GenerateTempl(ctx context.Context) (*dagger.Changeset, error) {
	ctr, err := m.buildCtr(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting build container: %w", err)
	}

	gen := withTemplGenerate(ctr).Directory(".")

	return m.Src.WithDirectory(".", gen).Changes(m.Src), nil
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
		return fmt.Errorf("templates are not up to date")
	}

	return nil
}
