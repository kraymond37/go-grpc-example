package main

import (
	"context"
	"github.com/kraymond37/go-grpc-example/proto"
)

type Backend struct {
}

func (b *Backend) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	reply := &proto.HelloReply{Message: req.Name}
	return reply, nil
}
