#!/bin/bash

set -o errexit

if [ ! -d "$(dirname "$1")" ] || [ ! -d "$(dirname "$2")" ] ; then
    echo "usage: $(basename "$0") <path-to-cmd-package> <path-to-output-file>" 1>&2
    exit 1
fi

CMD_PACKAGE_DIR=$1
EXECUTABLE_PATH=$2
APP_NAME=$3
GO_SRC_FILES=$(find "$CMD_PACKAGE_DIR" -name "*.go" | tr "\n" " ")

echo_call() {
    echo "$@"
    "$@"
}

# shellcheck disable=SC2086
echo_call go build -v \
    -o "$EXECUTABLE_PATH" \
    -ldflags="-X main.appID=$APP_NAME" \
    $GO_SRC_FILES
