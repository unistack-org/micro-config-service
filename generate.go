package service

//go:generate protoc -I./proto -I. --go-grpc_out=paths=source_relative:./proto --go_out=paths=source_relative:./proto --go-micro_out=components=micro|grpc,standalone=true,debug=true,paths=source_relative:./micro proto/service.proto
