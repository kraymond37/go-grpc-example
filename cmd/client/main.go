package main

import (
	"context"
	"github.com/kraymond37/go-grpc-example/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//conn, err := grpc.Dial("localhost:56789", grpc.WithInsecure(), grpc.WithBlock())
	conn, err := grpc.DialContext(ctx, "localhost:10000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()

	client := proto.NewExampleClient(conn)
	reply, err := client.SayHello(context.Background(), &proto.HelloRequest{Name: "you"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(reply.Message)
}
