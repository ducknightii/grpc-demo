PB_FILES=$(shell find pb/helloworld -name *.proto)

protoc:
	protoc -I ./pb \
	--go_out ./pb --go_opt paths=source_relative \
	--go-grpc_out ./pb --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./pb --grpc-gateway_opt paths=source_relative \
	$(PB_FILES)
