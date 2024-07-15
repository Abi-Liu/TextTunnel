#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd server/sql/schema
goose postgres $DATABASE_URL up
