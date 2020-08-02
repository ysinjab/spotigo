package main

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"

	pb "github.com/ysinjab/spotigo/pkg/albums"
	grpc "google.golang.org/grpc"
)

const (
	address = "localhost:8888"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAlbumsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetAlbum(ctx, &wrappers.Int32Value{Value: 1})
	if err != nil {
		log.Fatalf("could not get anything: %v", err)
	}
	log.Printf(r.Name)
}
