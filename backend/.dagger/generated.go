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
		Directory("genpb")

	return m.Src.WithDirectory("genpb", gen).Changes(m.Src)
}

// GenerateSqlc - generate sqlc codegen from .sql files
// +generate
func (m *Backend) GenerateSqlc() *dagger.Changeset {
	gensql := dag.Container().
		From("sqlc/sqlc:1.29.0").
		WithWorkdir("workdir").
		WithDirectory("sql", m.Src.Directory("sql")).
		WithFile("sqlc.yaml", m.Src.File("sqlc.yaml")).
		WithExec([]string{"generate"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Directory("gensql")

	return m.Src.WithDirectory("gensql", gensql).Changes(m.Src)
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

func withTemplGenerate(ctr *dagger.Container) *dagger.Container {
	return ctr.WithExec([]string{"go", "tool", "templ", "generate"})
}
