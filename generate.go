package service

//go:generate protoc -I./proto -I. --go_out=paths=source_relative:./proto proto/service.proto
//go:generate protoc -I./proto -I. --go-micro_out=components=micro,standalone=false,debug=true,paths=source_relative:./proto proto/service.proto
//go:generate protoc -I./proto -I. --go-micro_out=components=grpc,standalone=true,debug=true,paths=source_relative:./micro proto/service.proto
