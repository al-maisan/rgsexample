package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/al-maisan/rgsproto/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Check if client ID is provided
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <client_id>", os.Args[0])
	}
	clientID := os.Args[1]
	// Set up a connection to the server.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new client
	client := pb.NewEventServiceClient(conn)

	// Prepare the subscription request
	req := &pb.SubscribeRequest{
		ClientId: clientID, // You can generate a unique ID here
	}

	// Call the SubscribeEvents RPC
	stream, err := client.SubscribeEvents(ctx, req)
	if err != nil {
		log.Fatalf("Error on subscribe: %v", err)
	}

	// Continuously receive events
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			// Server has closed the stream
			log.Println("Server closed the stream")
			break
		}
		if err != nil {
			log.Fatalf("Error receiving event: %v", err)
		}

		// Process the received event
		fmt.Printf("Received event: ID=%s, Timestamp=%d, Data=%s\n",
			event.Id, event.Timestamp, event.Data)
	}
}
