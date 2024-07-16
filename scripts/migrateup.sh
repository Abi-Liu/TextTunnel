#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd internal/server/sql/schema
goose postgres $DATABASE_URL up
