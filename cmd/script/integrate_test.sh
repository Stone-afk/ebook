#!/usr/bin/env bash

set -e
docker compose -f cmd/docker-compose.yaml down
docker compose -f cmd/docker-compose.yaml up -d
go test -race ./... -tags=e2e
docker compose -f cmd/docker-compose.yaml down
