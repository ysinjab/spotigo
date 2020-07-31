package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ysinjab/spotigo/pkg/albums"
	pb "github.com/ysinjab/spotigo/pkg/albums"

	grpc "google.golang.org/grpc"
)

type AlbumsServer struct {
	pb.UnimplementedAlbumsServer
}

func (s *AlbumsServer) GetAlbums(context.Context, *pb.Empty) (*pb.Album, error) {
	return &albums.Album{Id: 1, Name: "D'You know what I mean"}, nil
}
func main() {
	fmt.Print("Heeeey")
	server := grpc.NewServer()
	albums.RegisterAlbumsServer(server, &AlbumsServer{})
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(":(")
	}
	log.Fatal(server.Serve(l))

}
