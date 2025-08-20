// A generated module for Backend functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/backend/internal/dagger"
)

type Backend struct{}

func (m *Backend) GenerateProtos(
	ctx context.Context,
	// +defaultPath="/backend"
	src *dagger.Directory,
) *dagger.Directory {

	return dag.Container().
		From("bufbuild/buf:latest").
		WithExec([]string{"apk", "add", "--no-cache", "protobuf-dev"}).
		WithWorkdir("workdir").
		WithDirectory("proto", src.Directory("proto")).
		WithFiles(".", []*dagger.File{src.File("buf.yaml"), src.File("buf.gen.yaml")}).
		WithExec([]string{"buf", "generate"}).
		Directory("gen")

}

// Returns lines that match a pattern in the files of the provided Directory
func (m *Backend) GrepDir(ctx context.Context, directoryArg *dagger.Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}
