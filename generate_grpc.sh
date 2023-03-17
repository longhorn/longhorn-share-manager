#!/bin/bash

set -e

# check and download dependency for gRPC code generate
if [ ! -e ./proto/vendor/protobuf/src/google/protobuf ]; then
    rm -rf ./proto/vendor/protobuf/src/google/protobuf
    DIR="./proto/vendor/protobuf/src/google/protobuf"
    mkdir -p $DIR
    wget https://raw.githubusercontent.com/protocolbuffers/protobuf/v3.9.0/src/google/protobuf/empty.proto -P $DIR
fi

# share manager rpc
python3 -m grpc_tools.protoc -I pkg/rpc -I proto/vendor/protobuf/src/ --python_out=integration/rpc/share_manager --grpc_python_out=integration/rpc/share_manager pkg/rpc/smrpc.proto
protoc -I pkg/rpc/ -I proto/vendor/protobuf/src/ pkg/rpc/smrpc.proto --go_out=plugins=grpc:pkg/rpc
