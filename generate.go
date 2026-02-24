package service

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v4) --go_out=paths=source_relative:./proto proto/service.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v4) --go-micro_out=components=micro,debug=true,paths=source_relative:./proto proto/service.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v4) --go-micro_out=components=grpc,standalone=true,debug=true,paths=source_relative:./grpc proto/service.proto"

//go:generate sh -c "protoc -I./proto -I. -I$(go list -f '{{ .Dir }}' -m go.unistack.org/micro-proto/v4) --go-micro_out=components=http,standalone=true,debug=true,paths=source_relative:./http proto/service.proto"
