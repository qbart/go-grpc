.PHONY: protoc
protoc:
	mkdir -p pb
	protoc --proto_path=. \
	  --go-grpc_opt=paths=source_relative --go-grpc_out=pb \
   	  --go_opt=paths=source_relative --go_out=pb port_domain_service.proto