package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/Qwilt/grpc-ingress/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func loadTLSCliCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	//pemServerCA, err := ioutil.ReadFile(<Private CA file)
	//if err != nil {
	//	return nil, err
	//}
	//
	certPool := x509.NewCertPool()
	//if !certPool.AppendCertsFromPEM(pemServerCA) {
	//	return nil, fmt.Errorf("failed to add server CA's certificate")
	//}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs:            certPool,
		ServerName:         "<my domain>",
		InsecureSkipVerify: true,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	tlsCredentials, err := loadTLSCliCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	var conn *grpc.ClientConn
	conn, err = grpc.Dial("<my domain>:443", grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := chat.NewChatServiceClient(conn)

	response, err := c.SayHello(context.Background(), &chat.Message{Body: "I'm the client"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %+v", response.Body)
}
