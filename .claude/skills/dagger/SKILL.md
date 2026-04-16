---
name: dagger
description: Edit and maintain dagger module code.
---

## What I do

- Edit and maintain dagger module code

## When to use me

Use this when you need to make code changes to dagger functions.

## What is Dagger

Dagger is a framework for CI/CD. Dagger modules are go projects that call the dagger API.

Dagger modules have two parts:
  - dagger codegen - lives in `dagger.gen.go` and `internal/dagger.gen.go`
  - user defined functions - lives in all other go files.
  
Do not read/edit dagger codegen files.

## Dagger modules

A `dagger.json` file defines a dagger module. Module source is located
in the "dagger module directory".

When examining `dagger.json` the "dagger module directory" is:

- The value of `.source`
- If `.source` doesn't exist, the directory containing `dagger.json`

### With `.source` example

`my-project/dagger.json`:

```json
{
  "name": "my-project",
  "engineVersion": "v0.20.3",
  "source": ".dagger"
}
```

The "dagger module directory" is: `my-project/.dagger`.

### Without `.source` example:

`my-project/dagger.json`:

```json
{
  "name": "my-project",
  "engineVersion": "v0.20.3"
}
```

The "dagger module directory" is: `my-project`.

## Dagger functions

Exported go functions are dagger functions:

```go
func (m *MainModule) GreetPerson(name string) error {
    fmt.Println("hello " + name)
}
```

Functions are called on the command line like so

```bash
dagger call greet-person --name=chris
```

Generalised as:

```bash
dagger call <function-name> --<param1>=value ... --<param-n>=value
```

Function names and parameters must be specified in kebab case when using them over the command line.

## Dagger function signatures

Dagger functions signatures can be in 1 of 3 formats:

1. Do not accept the context parameter, do not return error:

```go
func (m *MainModule) BuildContainer() *dagger.Container {
  return dag.Container().From("golang:latest")
}
```

2. Do not accept the context parameter, return error:

```go
func (m *MainModule) BuildContainer() (*dagger.Container, error) {
  return dag.Container().From("golang:latest"), nil
}
```

3. Accept the context parameter, do return an error:

```go
func (m *MainModule) BuildContainer(ctx context.Context) (*dagger.Container, error) {
  ver, err := goVersion(ctx, m.Src)
  if err != nil {
    return fmt.Errorf("error retrieving version from project source: %w", err)
  }

  return dag.Container().From("golang:" + ver)
}
```

The following is incorrect (accept the context, do not return error): 

```go
func (m *MainModule) BuildContainer(ctx context.Context) *dagger.Container {
 // ...
}
```

## Dagger checks

Dagger functions decorated with the `// +check` pragma are a "dagger check".
A "dagger check" does not accept parameters and returns an error:

```go
// +check
func (m *MainModule) CheckConfigFile(ctx context.Context) error {
    exists, err := m.Src.File("config.json").Exists(ctx)
    if err != nil {
        return fmt.Errorf("error occured when calling exists: %w", err)
    }

    if !exists {
        return errors.New("config.json file does not exist")
    }

    return nil
}
```

To run all dagger checks, use this command:

```bash
dagger checks
```

To list dagger checks, use this command:

```bash
dagger checks -l
```

## Dagger generate functions

Dagger functions decorated with the `// +generate` pragma are a "generate function".  
A "generate function" must not accept any parameters and must return a `*dagger.Changeset`.

```go
// +generate
func (m *MainModule) GenerateTestData(ctx context.Context) (*dagger.Changeset, error) {
    testsJson, err := generateTestData("tests.json", 3)
    if err != nil {
        return fmt.Errorf("failed to generate test data: %w", err)
    }

    gen := m.Src.WithNewFile("tests.json", testsJson)

    return gen.Changes(m.Src)
}
```

To run all dagger generate functions, use the this command:

```bash
dagger generate -y
```

To list generate functions, use the this command:

```bash
dagger generate -l
```

## Code changes workflow

Follow this workflow when making changes:

1. Determine the "dagger module directory" from the `dagger.json`.
2. Make changes to files in the "dagger module directory".
3. Run `dagger develop` from the "dagger module directory" to regenerate dagger codegen.
3. Run `go vet ./...` in the "dagger module directory", to verify there are no errors.
4. If `go vet ./...` passes, call `dagger checks`.