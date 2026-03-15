#!/bin/sh

if [ ! -n "${DATABASE_FILE}" ]; then
    echo "expected env var DATABASE_FILE to be set"
    exit 1
fi

if [ ! -f "${DATABASE_FILE}" ]; then
    echo "skipping migration as database doesn't exist"
    exit 0
fi

echo "running migration on database"

atlas schema apply \
  --url "sqlite:///${DATABASE_FILE}" \
  --to file:///app/new.sql \
  --dev-url "sqlite://dev?mode=memory"