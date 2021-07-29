.PHONY: test
test:
	go test ./...

.PHONY: protoc
protoc:
	mkdir -p pb
	protoc --proto_path=. \
	  --go-grpc_opt=paths=source_relative --go-grpc_out=pb \
   	  --go_opt=paths=source_relative --go_out=pb port_domain_service.proto

.PHONY: install_proto
install_proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

.PHONY: get
get:
	curl -i localhost:3000/ports/AEAJM

.PHONY: get_invalid
get_invalid:
	curl -i localhost:3000/ports/AEAJ
