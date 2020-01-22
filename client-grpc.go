package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/xshyne88/clientgrpc/helloworld"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:3001"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	fmt.Println("connected..")
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c := pb.NewStatGeneratorClient(conn)
	r, err := c.SendStats(ctx, &pb.Stats{})
	log.Printf("message received: \n memory: %v \n name: %s, \n os: %v", r.GetMemory(), r.GetName(), r.GetOs())

	fmt.Sprintf("c: %v %v", conn, c)
	select {}
}
