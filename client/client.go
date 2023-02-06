package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	_ "github.com/edwardOWO/goexample/learn"
	_ "github.com/edwardOWO/goexample/msg"
	pb "github.com/edwardOWO/goexample/route"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc"
)

func main() {

	StreamGrpcClient()
}

func SimpleGrpcClient() {

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := pb.NewGetTimeClient(conn)

	test := pb.TimeRequest{Timezone: 0}
	feature, err := client.GetTime(context.Background(), &test)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(feature)

}

func StreamGrpcClient() {

	rand.Seed(time.Now().Unix())

	// dail server
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	client := pb.NewStreamServiceClient(conn)
	in := &pb.StreamRequest{Id: 1}
	stream, err := client.FetchResponse(context.Background(), in)
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	//ctx := stream.Context()
	done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true //close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			log.Printf("Resp received: %s", resp.Result)
		}
	}()

	<-done
	log.Printf("finished")

}
