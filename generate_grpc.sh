#!/bin/bash

set -e

# share manager
python3 -m grpc_tools.protoc -I pkg/rpc --python_out=integration/rpc/share_manager --grpc_python_out=integration/rpc/share_manager pkg/rpc/rpc.proto
protoc -I pkg/rpc/ pkg/rpc/rpc.proto --go_out=plugins=grpc:pkg/rpc
