package service

//go:generate protoc -I./proto -I. --go-grpc_out=paths=source_relative:./proto --go_out=paths=source_relative:./proto --micro_out=components=micro|rpc,standalone=true,debug=true,paths=source_relative:./micro proto/service.proto
