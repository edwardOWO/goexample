// server.go
package server

import (
	"context"
	"fmt"
	"log"

	"time"

	pb "github.com/edwardOWO/goexample/route"
)

type SimpleGrpc struct {
	pb.UnimplementedGetTimeServer
}

/*
	func (s *routeGuideServer) GetFeature(cxt context.Context, point *pb.Point) (*pb.Feature, error) {
		for _, feature := range s.features {
			if proto.Equal(feature.Location, point) {
				return feature, nil
			}
		}
		return nil, nil
	}
*/
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

	time.Sleep(time.Duration(3) * time.Second)
	resp := pb.StreamResponse{Result: fmt.Sprintf("Request #%d For Id:%d", 1, in.Id)}
	if err := srv.Send(&resp); err != nil {
		log.Printf("send error %v", err)
	}
	time.Sleep(time.Duration(3) * time.Second)
	resp = pb.StreamResponse{Result: fmt.Sprintf("Request #%d For Id:%d", 1, in.Id)}
	if err := srv.Send(&resp); err != nil {
		log.Printf("send error %v", err)
	}
	time.Sleep(time.Duration(3) * time.Second)
	resp = pb.StreamResponse{Result: fmt.Sprintf("Request #%d For Id:%d", 1, in.Id)}
	if err := srv.Send(&resp); err != nil {
		log.Printf("send error %v", err)
	}
	time.Sleep(time.Duration(3) * time.Second)
	resp = pb.StreamResponse{Result: fmt.Sprintf("Request #%d For Id:%d", 1, in.Id)}
	if err := srv.Send(&resp); err != nil {
		log.Printf("send error %v", err)
	}

	return nil
}
