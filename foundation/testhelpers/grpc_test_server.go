package testhelpers

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func InitiateGRPCTestServer(interceptor grpc.UnaryServerInterceptor) *grpc.Server {
	lis = bufconn.Listen(bufSize)
	var server *grpc.Server
	if interceptor != nil {
		server = grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	} else {
		server = grpc.NewServer()
	}

	return server
}

func StartGRPCTestServer(server *grpc.Server) *grpc.Server {
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	return server
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func GetTestConnection() (*grpc.ClientConn, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
