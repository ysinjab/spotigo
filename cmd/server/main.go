package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/ysinjab/spotigo/pkg/albums"
	pb "github.com/ysinjab/spotigo/pkg/albums"
	"google.golang.org/grpc"
)

type albumsServer struct {
}

func (s *albumsServer) GetAlbums(context.Context, *pb.Empty) (*pb.Album, error) {
	return &albums.Album{Id: 1, Name: "D'You know what I mean"}, nil
}

func (s *albumsServer) GetAlbum(ctx context.Context, in *pb.AlbumId) (*pb.Album, error) {
	if in.Id == 1 {
		return &albums.Album{Id: 1, Name: "D'You know what I mean"}, nil
	}
	return nil, errors.New("Not found :(")
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
