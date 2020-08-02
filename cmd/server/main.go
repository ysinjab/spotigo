package main

import (
	"context"
	"log"
	"net"

	"github.com/ysinjab/spotigo/pkg/albums"
	pb "github.com/ysinjab/spotigo/pkg/albums"
	"google.golang.org/grpc"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"
)

type albumsServer struct {
}

func (s *albumsServer) GetAlbums(context.Context, *pb.Empty) (*pb.Album, error) {
	return &albums.Album{Id: 1, Name: "D'You know what I mean"}, nil
}

// GetAlbum(context.Context, *wrappers.Int32Value) (*Album, error)
func (s *albumsServer) GetAlbum(ctx context.Context, in *wrappers.Int32Value) (*pb.Album, error) {
	if in.GetValue() == 1 {
		return &albums.Album{Id: 1, Name: "D'You know what I mean"}, nil
	}
	return nil, nil
}

func main() {
	server := grpc.NewServer()
	albums.RegisterAlbumsServer(server, &albumsServer{})
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(":(")
	}
	log.Fatal(server.Serve(l))
}
