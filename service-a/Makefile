genpb:
	protoc --proto_path=adapter/service_b_adapter/pb adapter/service_b_adapter/pb/*.proto --go_out=adapter/service_b_adapter/pb --go_opt=paths=source_relative --go-grpc_out=adapter/service_b_adapter/pb --go-grpc_opt=paths=source_relative
	
start:
	go run cmd/*.go start

.PHONY: genpb start