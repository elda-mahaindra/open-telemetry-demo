genpb:
	protoc --proto_path=api/pb api/pb/*.proto --go_out=api/pb --go_opt=paths=source_relative --go-grpc_out=api/pb --go-grpc_opt=paths=source_relative
	
start:
	go run cmd/*.go start

.PHONY: start