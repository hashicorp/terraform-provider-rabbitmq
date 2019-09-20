#!/bin/bash
source .version
set -x

PROJECT_NAME=terraform-provider-rabbitmq
PKG=github.com/samueldumont/$PROJECT_NAME
VERSION=${MAJOR}.${MINOR}.${PATCH}

for arch in linux darwin; do
  GOFLAGS=-mod=vendor GO111MODULE=on GOOS=$arch GOARCH=amd64 go build -o "bin/$PROJECT_NAME-v$VERSION-$arch-amd64" .
done

GOFLAGS=-mod=vendor GO111MODULE=on CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o "bin/$PROJECT_NAME-v$VERSION-windows-amd64.exe" .
