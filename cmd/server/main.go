package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/al-maisan/rgsproto/api"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedEventServiceServer
	mu      sync.Mutex
	clients map[string]chan *pb.Event
}

func (s *server) SubscribeEvents(req *pb.SubscribeRequest, stream pb.EventService_SubscribeEventsServer) error {
	clientID := req.ClientId
	eventChan := make(chan *pb.Event, 100)

	s.mu.Lock()
	s.clients[clientID] = eventChan
	log.Printf("connected to client: '%s'", clientID)
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, clientID)
		close(eventChan)
		s.mu.Unlock()
	}()

	for {
		select {
		case event := <-eventChan:
			if err := stream.Send(event); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (s *server) broadcastEvent(event *pb.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, client := range s.clients {
		select {
		case client <- event:
		default:
			// If client's channel is full, skip this client
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	eventServer := &server{
		clients: make(map[string]chan *pb.Event),
	}
	pb.RegisterEventServiceServer(s, eventServer)

	go func() {
		for i := 0; ; i++ {
			event := &pb.Event{
				Id:        fmt.Sprintf("event-%d", i),
				Timestamp: time.Now().Unix(),
				Data:      fmt.Sprintf("Event data %d", i),
			}
			eventServer.broadcastEvent(event)
			time.Sleep(5 * time.Second)
		}
	}()

	log.Println("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
