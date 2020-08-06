#!/usr/bin/env make

define HELP_INFO=
Usage:
  make install	安装gRPC插件
  make generate	生成gRPC代码
  make server	编译gRPC server
  make client	编译gRPC client
  make clean	清理编译输出
endef
export HELP_INFO

.PHONY: help generate install server client clean

help:
	@echo "$$HELP_INFO"

generate:
	# Generate go, gRPC-Gateway, swagger, go-validators output.
	#
	# -I declares import folders, in order of importance
	# This is how proto resolves the protofile imports.
	# It will check for the protofile relative to each of these
	# folders and use the first one it finds.
	#
	# --go_out generates Go Protobuf output with gRPC plugin enabled.
	# --grpc-gateway_out generates gRPC-Gateway output.
	# --swagger_out generates an OpenAPI 2.0 specification for our gRPC-Gateway endpoints.
	# --govalidators_out generates Go validation files for our messages types, if specified.
	# --go-grpc_out generates Go language bindings of services in protobuf definition files for gRPC.
	#
	# ./proto is the output directory.
	#
	# proto/example.proto is the location of the protofile we use.
	protoc \
		-I proto \
		-I "${GOPATH}"/src/github.com/grpc-ecosystem/grpc-gateway/ \
		-I "${GOPATH}"/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/ \
		-I "${GOPATH}"/src/ \
		--go_out=plugins=grpc,paths=source_relative:./proto/ \
		--grpc-gateway_out=allow_patch_feature=false,paths=source_relative:./proto/ \
		--swagger_out=swagger/ \
		--govalidators_out=paths=source_relative:./proto \
		proto/*.proto

	# Generate static assets for OpenAPI UI
	statik -m -f -src swagger/

install:
	# install protobuf
	# https://github.com/protocolbuffers/protobuf/blob/master/src/README.md

	# install plugins
	go get \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/mwitkow/go-proto-validators/protoc-gen-govalidators \
		github.com/rakyll/statik

server:
	go build -v ./cmd/server

client:
	go build -v ./cmd/client

clean:
	go clean --cache
	$(RM) server client