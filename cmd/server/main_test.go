package main_test

import (
	context "context"
	"fmt"
	"log"
	"net"
	"testing"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	main "github.com/ysinjab/spotigo/cmd/server"
	"github.com/ysinjab/spotigo/pkg/albums"
	pb "github.com/ysinjab/spotigo/pkg/albums"
	"github.com/ysinjab/spotigo/pkg/auth"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

// type testingAlbumsServer struct {
// 	service albums.Service
// }

// func (s *testingAlbumsServer) GetAlbums(context.Context, *pb.Empty) (*pb.AlbumList, error) {
// 	return &pb.AlbumList{Albums: []*albums.Album{}}, nil
// }

// func (s *testingAlbumsServer) GetAlbum(ctx context.Context, in *pb.AlbumId) (*pb.Album, error) {
// 	return &albums.Album{Id: 1, Name: "D'You know what I mean"}, nil
// }

type AlbumsTestSuite struct {
	suite.Suite

	AlbumsService pb.AlbumsServer
	ServerOpts    []grpc.ServerOption
	ClientOpts    []grpc.DialOption

	serverAddr     string
	ServerListener net.Listener
	Server         *grpc.Server
	clientConn     *grpc.ClientConn
	Client         pb.AlbumsClient
}
type repository struct {
}

func (s *repository) GetAlbums() ([]albums.Album, error) {
	list := []albums.Album{}
	return list, nil
}

func (suite *AlbumsTestSuite) SetupTest() {
	lis = bufconn.Listen(bufSize)
	suite.Server = grpc.NewServer(suite.ServerOpts...)

	albumsService, err := albums.NewService(&repository{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	ss := &main.AlbumsServer{Service: albumsService}

	pb.RegisterAlbumsServer(suite.Server, ss)
	go func() {
		if err := suite.Server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func ctxWithToken(ctx context.Context, scheme string, token string) context.Context {
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %v", scheme, token))
	nCtx := metautils.NiceMD(md).ToOutgoing(ctx)
	return nCtx
}

func (s *AlbumsTestSuite) TestGetAlbums() {
	ctx := ctxWithToken(context.Background(), "Bearer", "eyJhbGciOiJIUzI1NiIsIR5cCI6IkpXVCJ9.eyJleHAiOjE1OTg4MjU3NjIsInN1YiI6MTIzfQ.5xs5nh3n-9wAbCKeyWM6JZJYZBtxP07O17kqp64oBmM")
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		s.Error(err)
	}
	defer conn.Close()

	client := pb.NewAlbumsClient(conn)
	r, err := client.GetAlbum(ctx, &pb.AlbumId{Id: 1})
	if err != nil {
		s.Error(err)
	}

	require.NoError(s.T(), err, "should not fail on establishing the stream")
	s.NotNil(r)
}

type AuthTestSuite struct {
	suite.Suite
}

func TestAuth(t *testing.T) {
	// authFunc := buildDummyAuthFunction("bearer", commonAuthToken)
	s := &AlbumsTestSuite{
		ServerOpts: []grpc.ServerOption{
			grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(auth.AuthFunc)),
			grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(auth.AuthFunc)),
		},
	}
	suite.Run(t, s)
}
