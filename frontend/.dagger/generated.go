package main

import (
	"context"
	"dagger/frontend/internal/dagger"
	"fmt"
)

// GenerateProtos - generate protobuf codegen from .proto files
// +generate
func (m *Frontend) GenerateProtos(ctx context.Context) (*dagger.Changeset, error) {
	nodeVer, err := m.nodeVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting node version: %w", err)
	}

	gen := m.buildCtr(nodeVer).
		WithExec([]string{"npm", "install"}).
		WithExec([]string{"npm", "run", "generate-protos"}).
		Directory("src/gen")

	return m.Src.WithDirectory("src/gen", gen).Changes(m.Src), nil
}

// CheckProtos - check that the working tree's proto generated files are in sync.
// +check
func (m *Frontend) CheckProtos(ctx context.Context) error {
	chgset, err := m.GenerateProtos(ctx)
	if err != nil {
		return fmt.Errorf("error generating protos: %w", err)
	}

	return assertEmpty(ctx, chgset)
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
