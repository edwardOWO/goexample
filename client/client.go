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

	ServerSide()
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

// Client Side Streaming
func ClientSide() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立gRPC连接
	grpcClient := pb.NewClientSideClient(conn)
	// 创建发送结构体
	res, err := grpcClient.ClientSideHello(context.Background())
	if err != nil {
		log.Fatalf("Call SayHello err: %v", err)
	}

	for i := 0; i < 5; i++ {
		//通过 Send方法发送流信息
		err = res.Send(&pb.ClientSideRequest{Id: 123})
		if err != nil {
			return
		}
	}

	// 打印返回值
	log.Println(res.CloseAndRecv())
}

func ServerSide() {
	// 连接服务器
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立gRPC连接
	grpcClient := pb.NewServerSideClient(conn)
	// 创建发送结构体
	req := pb.ServerSideRequest{
		Id: 1,
	}
	//获取流
	stream, err := grpcClient.ServerSideHello(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call SayHello err: %v", err)
	}
	for n := 0; n < 5; n++ {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Fatalf("No data: %v", err)
			break
		}
		if err != nil {
			log.Fatalf("Conversations get stream err: %v", err)
		}
		// 打印返回值
		log.Println(res.Result)
	}
}
