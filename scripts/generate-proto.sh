# Generates the proto files for the action plugin - needs to be run from the root of the repo

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative plugins/action.proto
