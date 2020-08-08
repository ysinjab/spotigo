package main

import (
	"context"
	"fmt"
	"log"
	"time"

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
	r, err := c.GetAlbums(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get anything: %v", err)
	}
	s := fmt.Sprintf("Total returned: %d", len(r.Albums))
	log.Printf(s)
}
