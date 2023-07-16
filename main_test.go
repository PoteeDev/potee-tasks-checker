package main

import (
	"context"
	"log"
	"sync"
	"testing"

	pb "github.com/PoteeDev/potee-tasks-checker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	services  = []string{"test", "test", "test", "test", "test"}
	hostCount = 10
)

func ParceReply(reply *pb.Reply) {
	status := 0
	for _, row := range reply.Results {
		status += int(row.Status)
	}
	log.Println(status)
}

func TestClientServer(t *testing.T) {
	conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Code removed for brevity

	client := pb.NewCheckerClient(conn)

	// Note how we are calling the GetBookList method on the server
	// This is available to us through the auto-generated code
	wg := &sync.WaitGroup{}

	for _, service := range services {
		svc := service
		wg.Add(1)
		go func() {
			// Generate Ping data
			pingReq := &pb.PingRequest{
				Service: svc,
			}
			for i := 1; i < hostCount+1; i++ {
				pingReq.Data = append(pingReq.Data, &pb.PingData{
					Host: "localhost",
					Id:   int64(i),
				})
			}
			pingReply, _ := client.Ping(context.Background(), pingReq)
			log.Println(svc, "ping: ")
			ParceReply(pingReply)

			// Generate Get data
			getReq := &pb.GetRequest{
				Service: svc,
				Name:    "example",
			}
			for i := 1; i < hostCount+1; i++ {
				getReq.Data = append(getReq.Data, &pb.GetData{
					Host:  "localhost",
					Id:    int64(i),
					Value: "1",
				})
			}
			getReply, _ := client.Get(context.Background(), getReq)
			log.Println(svc, "get: ")
			ParceReply(getReply)

			// Generate Put data
			putReq := &pb.PutRequest{
				Service: svc,
				Name:    "example",
			}
			for i := 1; i < hostCount+1; i++ {
				putReq.Data = append(putReq.Data, &pb.PutData{
					Host: "localhost",
					Id:   int64(i),
					Flag: "qwe",
				})
			}

			putReply, _ := client.Put(context.Background(), putReq)
			log.Println(svc, "put: ")
			ParceReply(putReply)

			// Generate Exploit data
			exploitReq := &pb.ExploitRequest{
				Service: svc,
				Name:    "rce",
			}
			for i := 1; i < hostCount+1; i++ {
				exploitReq.Data = append(exploitReq.Data, &pb.ExploitData{
					Host: "localhost",
					Id:   int64(i),
				})
			}
			exploitReply, _ := client.Exploit(context.Background(), exploitReq)
			log.Println(svc, "exploit: ")
			ParceReply(exploitReply)

			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkTestServer(t *testing.B) {
	conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Code removed for brevity

	client := pb.NewCheckerClient(conn)

	// Note how we are calling the GetBookList method on the server
	// This is available to us through the auto-generated code
	reply, _ := client.Ping(context.Background(), &pb.PingRequest{})
	log.Println(reply)
}
