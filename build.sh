#!/bin/bash
source .version
set -x

PROJECT_NAME=terraform-provider-rabbitmq
PKG=github.com/samueldumont/$PROJECT_NAME
VERSION=${MAJOR}.${MINOR}.${PATCH}

for arch in linux darwin; do
  GOFLAGS=-mod=vendor GO111MODULE=on GOOS=$arch GOARCH=amd64 go build -o "bin/$PROJECT_NAME" .
  zip -r -j "bin/${PROJECT_NAME}_v${VERSION}_${arch}_amd64.zip" bin/${PROJECT_NAME}
done

GOFLAGS=-mod=vendor GO111MODULE=on CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o "bin/$PROJECT_NAME" .
zip -r -j "bin/${PROJECT_NAME}_v${VERSION}_windows_amd64.zip" bin/${PROJECT_NAME}
