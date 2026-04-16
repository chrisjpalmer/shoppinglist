---
name: codebase
description: General information about the codebase, load me everytime
---

# About the codebase

IMPORTANT: Dont update generated code directly. Use `dagger generate -y` if you need to change generated code.

## Backend

- Directory: `/backend`
- Notes:
  - To regenerate code: `dagger generate -y`
  - To run all checks: `dagger checks`

## Frontend

- Directory: `/frontend`
- Notes:
  - To regenerate code: `dagger generate -y`
  - To run all checks: `dagger checks`

## Dagger code

- Directories: `**/.dagger`
- Load the dagger skill when working with these.

## Making commits

A commit message should contain a subject line stating what is changing.
It should contain a body that restates the subject line but adds the
reason *why* the code is changing. Avoid explaining how the code is changing.

When writing commit messages, use the gitmoji skill to pick a gitmoji that
is appropriate for the change occuring.

Here is an example:

```
🚚 Rename generated protobuf code
    
Rename generated protobuf code from `gen` to `genpb`
to align with the approach taken to generated sqlc code
as well as provide better clarity.
```