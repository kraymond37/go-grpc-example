package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kraymond37/go-grpc-example/proto"
	_ "github.com/kraymond37/go-grpc-example/statik"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"os"
)

var (
	log         grpclog.LoggerV2
	gRpcPort    = "10000"
	gatewayPort = "20000"
)

func init() {
	log = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)
}

func startRpcServer() {
	gRpcAddr := fmt.Sprintf("localhost:%v", gRpcPort)
	lis, err := net.Listen("tcp", gRpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gRpcServer := grpc.NewServer()
	proto.RegisterExampleServer(gRpcServer, &Backend{})
	reflection.Register(gRpcServer)

	log.Info("Serving gRPC on http://", gRpcAddr)
	if err := gRpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func startRpcGateway() {
	gRpcAddr := fmt.Sprintf("localhost:%v", gRpcPort)
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", gRpcAddr)
	gwMux := runtime.NewServeMux()
	err := proto.RegisterExampleHandlerFromEndpoint(context.Background(), gwMux, dialAddr, []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()})
	if err != nil {
		log.Fatalln("failed to register handler", err)
	}
	gatewayAddr := fmt.Sprintf("localhost:%v", gatewayPort)
	log.Fatalln(http.ListenAndServe(gatewayAddr, gwMux))
}

func startSwagger() {
	gRpcAddr := fmt.Sprintf("localhost:%v", gRpcPort)
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", gRpcAddr)
	gwMux := runtime.NewServeMux()
	err := proto.RegisterExampleHandlerFromEndpoint(context.Background(), gwMux, dialAddr, []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()})
	if err != nil {
		log.Fatalln("failed to register handler", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwMux)

	mime.AddExtensionType(".svg", "image/svg+xml")
	statikFS, err := fs.New()
	if err != nil {
		log.Fatalln(err)
	}
	// Expose files in static on <host>/openapi-ui
	fileServer := http.FileServer(statikFS)
	prefix := "/swagger/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))

	gatewayAddr := fmt.Sprintf("localhost:%v", gatewayPort)
	log.Info("Serving gRPC-Gateway on http://", gatewayAddr)
	log.Info("Serving OpenAPI Documentation on http://", gatewayAddr, prefix)
	log.Fatalln(http.ListenAndServe(":20000", mux))
}

func main() {
	go startRpcServer()
	//go startRpcGateway()
	go startSwagger()
	select {}
}
