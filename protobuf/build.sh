#!/bin/bash

#Generate protobuf implementation
protoc -I/usr/local/include -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. *.proto

#Generate resetful <-> grpc gateway implementation
protoc -I/usr/local/include -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. *.proto

#Generate swagger definition
protoc -I/usr/local/include -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. *.proto