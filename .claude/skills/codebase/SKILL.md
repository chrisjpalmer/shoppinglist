---
name: codebase
description: General information about the codebase, load me everytime
---

# About the codebase

IMPORTANT: Dont update generated code directly. Use `dagger generate -y` if you need to change generated code.

## Running all checks

- Directory: `/`
- Command: `dagger checks`

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

## Notes on generating code

Always run `dagger generate -y` from either the `frontend` or `backend` directory, not from the repo root.

**Why:** Running from the root generates files in the wrong locations (stray directories like genpb/, gensql/, shopping/, src/ appear at the repo root).

**How to apply:** `cd backend && dagger generate -y` or `cd frontend && dagger generate -y` depending on what changed.