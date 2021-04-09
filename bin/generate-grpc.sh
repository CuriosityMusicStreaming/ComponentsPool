#!/usr/bin/env bash

set -o errexit

GOPATH=$(go env GOPATH)
PATH=${GOPATH}/bin:${PATH}

# Prints path to subdirectory if it exist, otherwise prints nothing
find_subdir_safely() {
    local DIR=$1
    local NAME_PATTERN=$2
    if [ -d "$DIR" ]; then
        find "$DIR" -name "$NAME_PATTERN" | head -n 1
    fi
}

# Prints path to grpc-gateway source code directory
find_grpc_gateway_dir() {
    local GO_GET_PATH=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/
    local GO_MOD_PATH
    GO_MOD_PATH=$(find_subdir_safely "${GOPATH}/pkg/mod/github.com/grpc-ecosystem/" "grpc-gateway*")
    if [ -d "$GO_MOD_PATH" ]; then
        echo "$GO_MOD_PATH"
    elif [ -d "$GO_GET_PATH" ]; then
        echo "$GO_GET_PATH"
    else
        echo 'error: cannot find grpc-gateway dir, install it with go-get' 1>&2
        exit 1
    fi
}

# Prints and calls given command
echo_call() {
    echo "$@"
    "$@"
}

GPRC_GATEWAY_DIR=$(find_grpc_gateway_dir)

if [ -z "$1" ]; then
    echo "usage: $(basename "$0") <path-to-proto-file>" 1>&2
    exit 1
elif [ ! -f "$1" ]; then
    echo "file '$1' not exist" 1>&2
    exit 1
fi

PROTO_DIR=$(dirname "$1")
PROTO_NAME=$(basename "$1")
echo_call protoc \
    "-I/usr/include" \
    "-I${GOPATH}/src" \
    "-I${GPRC_GATEWAY_DIR}/third_party/googleapis" \
    "-I${PROTO_DIR}" \
    "--go_out=plugins=grpc:${PROTO_DIR}" \
    "--grpc-gateway_out=logtostderr=true:${PROTO_DIR}" \
    "--swagger_out=logtostderr=true:${PROTO_DIR}" \
    "${PROTO_DIR}/${PROTO_NAME}"
