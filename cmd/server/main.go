package main

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ysinjab/spotigo/pkg/albums"
	"github.com/ysinjab/spotigo/pkg/auth"
	"github.com/ysinjab/spotigo/pkg/storage/postgres"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	pb "github.com/ysinjab/spotigo/pkg/albums"
	"google.golang.org/grpc"
)

var publicMethods = []string{"/albums.Albums/GetAlbum"}

const secret = "123"

type AlbumsServer struct {
	Service albums.Service
}

func (s *AlbumsServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	for _, v := range publicMethods {
		if v == fullMethodName {
			log.Print("Yahoo !")
			return ctx, nil
		}
	}
	ctx, err := auth.ValidateToken(ctx)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (s *AlbumsServer) GetAlbums(context.Context, *pb.Empty) (*pb.AlbumList, error) {
	as, err := s.Service.GetAlbums()
	if err != nil {
		return nil, err
	}
	allAlbums := []*albums.Album{}
	for _, a := range as {
		allAlbums = append(allAlbums, &a)
	}
	return &pb.AlbumList{Albums: allAlbums}, nil
}

func (s *AlbumsServer) GetAlbum(ctx context.Context, in *pb.AlbumId) (*pb.Album, error) {
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
	log.Println("Gonna start the server .... 🚀🚀🚀🚀🚀")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": 123,
		"exp": float64(time.Now().Unix() + 100000),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal("😭")
	}
	log.Println("Enjoy ... your token is: ", tokenString)

	postgresStorage, err := postgres.NewStorage()
	if err != nil {
		log.Fatal("😭")
	}

	albumsService, err := albums.NewService(postgresStorage)
	if err != nil {
		log.Fatal("😭")
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(auth.AuthFunc)))
	// server := grpc.NewServer(grpc.UnaryInterceptor(getAlbumUnaryServerInterceptor))

	albums.RegisterAlbumsServer(server, &AlbumsServer{Service: albumsService})
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("😭")
	}
	log.Println("Listening at: localhost:8888")
	log.Fatal(server.Serve(l))
}
