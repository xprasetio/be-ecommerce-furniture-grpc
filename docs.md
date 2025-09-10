Running Protoc to generate

protoc --go_out=./pb --go-grpc_out=./pb --proto_path=./proto --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative service/service.proto
