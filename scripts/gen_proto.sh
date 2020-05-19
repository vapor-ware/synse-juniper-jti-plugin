#!/bin/bash

#
# gen_proto.sh
#
# A script to generate the compiled .proto files defined in
# pkg/protocol/jti/protos. This script should be run from the
# project root.
#

if ! [[ -x "$(command -v protoc)" ]]; then
    echo 'error: protoc is not installed' >&2
    echo -e '\nget the latest protoc release: https://github.com/protocolbuffers/protobuf/releases'
    echo "or, on macOS with Homebrew: 'brew install protoc-gen-go'"
    exit 1
fi

if ! [[ -x "$(command -v protoc-gen-go)" ]]; then
    echo 'error: protoc-gen-go is not installed' >&2
    echo -e '\ninstall with "go install google.golang.org/protobuf/cmd/protoc-gen-go"'
    exit 1
fi

proto_dir="$(pwd)/pkg/protocol/jti/"

if [[ ! -d "${proto_dir}" ]]; then
    echo "error: output directory not found (expected ${proto_dir})" >&2
    echo -e "\nensure you are running script from within project root"
    echo "e.g. './scripts/gen_proto.sh'"
    exit 1
fi

if [[ ! -d "./protos" ]]; then
    echo "error: 'protos' directory not found"  >&2
    echo -e "\nensure you are running the script from within the project root"
    echo "e.g. './scripts/gen_proto.sh'"
    exit 1
fi


echo "Compiling proto files in ${proto_dir}..."

count=0
for f in ./protos/*; do
    filename=${f##*/}
    name=${filename%.*}

    echo "• compiling ${name}"

    protoc --proto_path=./protos \
           --go_out="${proto_dir}" \
           "${f}" > /dev/null 2>&1  # Note - can comment this redirect out if compile is failing

    if [[ $? -ne 0 ]]; then
        echo "error: failed to compile .proto file ${name} - terminating"
        exit 1
    fi
    ((count++))
done

echo "successfully compiled ${count} proto definitions"
