package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ysinjab/spotigo/pkg/albums"
	"github.com/ysinjab/spotigo/pkg/storage/postgres"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	pb "github.com/ysinjab/spotigo/pkg/albums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type key int

const (
	secret           = "123"
	tokenInfoKey key = iota
)

var publicMethods = []string{"/albums.Albums/GetAlbum"}

func parseToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if parsedToken.Valid {
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			fmt.Println(claims)

			return claims, nil
		} else {
			return nil, fmt.Errorf("invalid toke")
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, fmt.Errorf("that's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, fmt.Errorf("token is either expired or not active yet")
		} else {
			return nil, fmt.Errorf("couldn't handle this token %s", err)
		}
	} else {
		return nil, fmt.Errorf("couldn't handle this token %s", err)
	}
}

func userClaimFromToken(claims jwt.MapClaims) interface{} {
	return claims["sub"]
}

func doTheAuth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	grpc_ctxtags.Extract(ctx).Set("auth.sub", userClaimFromToken(tokenInfo))
	newCtx := context.WithValue(ctx, tokenInfoKey, tokenInfo)
	return newCtx, nil
}

func authFunc(ctx context.Context) (context.Context, error) {
	newCtx, err := doTheAuth(ctx)
	if err != nil {
		return nil, err
	}
	return newCtx, nil
}

func (s *albumsServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	for _, v := range publicMethods {
		if v == fullMethodName {
			log.Print("Yahoo !")
			return ctx, nil
		}
	}
	ctx, err := doTheAuth(ctx)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

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
	var albumsService albums.Service
	albumsService, err = albums.NewService(postgresStorage)
	if err != nil {
		log.Fatal("😭")
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(authFunc)))
	// server := grpc.NewServer(grpc.UnaryInterceptor(getAlbumUnaryServerInterceptor))

	albums.RegisterAlbumsServer(server, &albumsServer{service: albumsService})
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("😭")
	}
	log.Println("Listening at: localhost:8888")
	log.Fatal(server.Serve(l))
}
