package main

import (
	"context"

	"fmt"
	"github.com/Qwilt/grpc-ingress/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"net/http"

	"log"
	"net"
)

const (
	grpcPort = ":8443"
	httpPort = ":8080"
)

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	return handler(ctx, req)
}
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func main() {

	http.HandleFunc("/", HealthCheck)
	go func() {
		log.Printf("starting http server on  %v", httpPort)
		http.ListenAndServe(httpPort, nil)
	}()

	creds, err := credentials.NewServerTLSFromFile("cert/cert.pem", "cert/key.pem")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	reflection.Register(grpcServer)
	s := chat.Server{}
	chat.RegisterChatServiceServer(grpcServer, &s)
	log.Printf("starting up GRPC server on  %v", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

}
