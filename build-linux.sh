#!/bin/bash -e

ORG_PATH=github.com/nildev
REPO_PATH="${ORG_PATH}/watcher"
VER=$1
if [ -z "$1" ]; then
    VER=`git rev-parse --abbrev-ref HEAD`
fi

echo "Building watcher"
GOOS=linux GOARCH=amd64 godep go build -o bin/watcher -a -installsuffix netgo -ldflags "-s -X main.Version=$VER -X main.GitHash=`git rev-parse HEAD` -X main.BuiltTimestamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" ${REPO_PATH}
chmod +x bin/watcher