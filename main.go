package main

import (
	"log"

	"github.com/TOMMy-Net/tages/internal/grpc"
)

func main() {
	server := grpc.NewServer()
	err := server.Serve("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
}
