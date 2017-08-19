#!/bin/sh
# https://golang.org/doc/install/source#environment
#

go version > /dev/null || { echo "Please install the Go toolchain, e.g. $ apt-get install golang-go" 1>&2 && exit 1; }

cd "$(dirname "${0}")"
# $ uname -s -m
# Darwin x86_64
# Linux x86_64
# Linux armv6l

PROG_NAME="denkmaeler-xml2ttl"
VERSION="0.0.2"

rm "${PROG_NAME}"-*-"${VERSION}" 2>/dev/null

cd "${PROG_NAME}-cmd"

# go get -u "github.com/stretchr/testify"

go fmt && go test -i

env GOOS=linux GOARCH=amd64 go build -o "${PROG_NAME}-Linux-x86_64-${VERSION}"
# http://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5
bin="${PROG_NAME}-$(uname -s)-$(uname -m)-${VERSION}"
go build -o "${bin}"

ls -Al "${bin}"
