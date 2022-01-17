PB_FILES=$(shell find pb -name *.proto)

protoc:
	protoc --go-grpc_out=paths=source_relative:. $(PB_FILES)
	protoc --go_out=paths=source_relative:.  $(PB_FILES)
