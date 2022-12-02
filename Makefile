# generates proto buff types and interfaces from .proto file
proto:
	protoc --go_out=./ --go_opt=paths=source_relative \
                    --go-grpc_out=./ --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
                    internal/service/**/proto/*.proto