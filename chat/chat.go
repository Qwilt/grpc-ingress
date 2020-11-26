package chat

import (
	"context"
	"log"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("in SayHello input %+v", in.Body)
	return &Message{Body: "Hello " + in.Body}, nil
}
