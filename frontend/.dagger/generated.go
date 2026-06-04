package main

import (
	"dagger/frontend/internal/dagger"
)

// GenerateProtos - generate protobuf codegen from .proto files
// +generate
func (m *Frontend) GenerateProtos() *dagger.Changeset {
	gen := m.buildCtr().
		WithExec([]string{"npm", "install"}).
		WithExec([]string{"npm", "run", "generate-protos"}).
		Directory("src/gen")

	return m.Src.WithDirectory("src/gen", gen).Changes(m.Src)
}
