#!/bin/bash

# Find all .proto files under the protocols directory and run protoc for each
find . -type f -name "*.proto" | while read -r proto_file; do
  # Get the directory of the .proto file
  dir=$(dirname "$proto_file")
  
  # Run the protoc command in the directory of the .proto file
  (
    cd "$dir" || exit
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative "$(basename "$proto_file")"
  )
done
