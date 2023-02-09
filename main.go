package main

import (
	"log"
	"net"

	_ "github.com/edwardOWO/goexample/learn"
	_ "github.com/edwardOWO/goexample/msg"
	"github.com/edwardOWO/goexample/server"
	_ "github.com/edwardOWO/goexample/server"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc"

	pb "github.com/edwardOWO/goexample/route"

	_ "github.com/edwardOWO/goexample/route"
)

func main() {

	// 生成一個listener
	/*
		listener, err := net.Listen("tcp", "localhost:5000")
		if err != nil {
			log.Fatalln("cannot create a listener a the address")
		}
		// server
		grpcServer := grpc.NewServer()
		pb.RegisterRouteGuideServer(grpcServer, server.DbServer())

		log.Fatalln(grpcServer.Serve(listener))
	*/
	println("gRPC server tutorial in Go")

	// simple
	/*
		listener, err := net.Listen("tcp", ":9000")
		if err != nil {
			panic(err)
		}

		s := grpc.NewServer()

		pb.RegisterGetTimeServer(s, &server.Server{})
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	*/

	// stream rpc

	/*
		listener, err := net.Listen("tcp", ":9000")
		if err != nil {
			panic(err)
		}

		s := grpc.NewServer()

		pb.RegisterStreamServiceServer(s, server.StreamGRPC{})
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	*/

	ServerSideStreamGrpc()

}
func ClientSideStreamGrpc() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	pb.RegisterClientSideServer(s, &server.ClientSideService{})

	//pb.RegisterStreamClientServer(s, server.StreamClient{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
func ServerSideStreamGrpc() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	pb.RegisterServerSideServer(s, &server.ServerSideService{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
