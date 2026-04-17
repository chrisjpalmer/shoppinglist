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