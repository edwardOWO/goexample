// server.go
package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/edwardOWO/goexample/route"
)

type SimpleGrpc struct {
	pb.UnimplementedGetTimeServer
}

// Simple mode
func (s *SimpleGrpc) GetTime(cxt context.Context, rq *pb.TimeRequest) (*pb.TimeReply, error) {

	t := time.Now().UTC().Unix()

	timeReply := pb.TimeReply{Timestamp: t}

	return &timeReply, nil
}

type StreamGRPC struct {
	pb.UnimplementedStreamServiceServer
}

func (s StreamGRPC) FetchResponse(in *pb.StreamRequest, srv pb.StreamService_FetchResponseServer) error {
	log.Printf("fetch response for id : %d", in.Id)
	/*
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(count int64) {
				defer wg.Done()
				time.Sleep(time.Duration(count) * time.Second)
				resp := pb.StreamResponse{Result: fmt.Sprintf("Request #%d For Id:%d", count, in.Id)}
				if err := srv.Send(&resp); err != nil {
					log.Printf("send error %v", err)
				}
				log.Printf("finishing request number : %d", count)
			}(int64(i))
		}

		wg.Wait()
	*/
	i := 0
	for {

		time.Sleep(time.Second * 1)
		i += 1
		resp := pb.StreamResponse{Result: fmt.Sprintf("Request #%d For Id:%d", i, in.Id)}

		if err := srv.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
	}

	return nil
}

type ClientSideService struct {
	pb.UnimplementedClientSideServer
}

// Client Side Streaming
func (c *ClientSideService) ClientSideHello(server pb.ClientSide_ClientSideHelloServer) error {
	for i := 0; i < 5; i++ {
		recv, err := server.Recv()
		if err != nil {
			return err
		} else if err == io.EOF {
			log.Println("No data")
			return err
		}
		log.Println("Client Message", recv)
	}
	err := server.SendAndClose(&pb.ClientSideResp{Result: "Close"})
	if err != nil {
		return err
	}
	return nil

}

// ServerSideStream
type ServerSideService struct {
	pb.UnimplementedServerSideServer
}

func (c *ServerSideService) ServerSideHello(request *pb.ServerSideRequest, server pb.ServerSide_ServerSideHelloServer) error {

	log.Println(request.Id)
	for n := 0; n < 3; n++ {
		time.Sleep(time.Second * 1)
		err := server.Send(&pb.ServerSideResp{Result: "Hello"})
		if err != nil {
			return err
		}
	}
	return nil
}
