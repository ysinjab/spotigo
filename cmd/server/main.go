package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/ysinjab/spotigo/pkg/albums"
	"github.com/ysinjab/spotigo/pkg/storage/postgres"

	pb "github.com/ysinjab/spotigo/pkg/albums"
	"google.golang.org/grpc"
)

type albumsServer struct {
	service albums.Service
}

func (s *albumsServer) GetAlbums(context.Context, *pb.Empty) (*pb.AlbumList, error) {
	as, err := s.service.GetAlbums()
	if err != nil {
		return nil, err
	}
	allAlbums := []*albums.Album{}
	for _, a := range as {
		allAlbums = append(allAlbums, &a)
	}
	return &pb.AlbumList{Albums: allAlbums}, nil
}

func (s *albumsServer) GetAlbum(ctx context.Context, in *pb.AlbumId) (*pb.Album, error) {
	if in.Id == 1 {
		return &albums.Album{Id: 1, Name: "D'You know what I mean"}, nil
	}
	return nil, errors.New("Not found :(")
}

func getAlbumUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("======= [Server Interceptor] ", info.FullMethod)
	m, err := handler(ctx, req)
	log.Printf(" Post Proc Message : %s", m)
	return m, err
}

func main() {
	// interceptor as an example
	postgresStorage, err := postgres.NewStorage()
	if err != nil {
		log.Fatal("T_T")
	}
	var albumsService albums.Service
	albumsService, err = albums.NewService(postgresStorage)
	if err != nil {
		log.Fatal("T_T")
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(getAlbumUnaryServerInterceptor))
	albums.RegisterAlbumsServer(server, &albumsServer{service: albumsService})
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(":(")
	}
	log.Fatal(server.Serve(l))
}
