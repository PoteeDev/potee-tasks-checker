package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/PoteeDev/potee-tasks-checker/proto"
	"google.golang.org/grpc"
)

func GetSrtingEnvDefault(env, defaults string) string {
	value := os.Getenv(env)
	if value != "" {
		return value
	}
	return defaults
}

func GetIntEnvDefault(env string, defaults int) int {
	value := os.Getenv(env)
	if value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Println(fmt.Sprintf("value %s of %s is not int, use %d", env, value, defaults))
			return defaults
		}
		return intValue
	}
	return defaults
}

var (
	port      = flag.Int("port", 50051, "The server port")
	timeout   = flag.Int("timeout", GetIntEnvDefault("TIMEOUT", 5), "timeot for request")
	directory = flag.String("dir", GetSrtingEnvDefault("DIR", "examples"), "directory with scripts")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedCheckerServer
}

func RunCheckers(command string, result chan []*pb.Result) {
	output, err := Execute(command)
	if err != nil {
		log.Println(err)
	}
	resultStruct := []*pb.Result{}

	log.Println(output)
	err = json.Unmarshal([]byte(output), &resultStruct)
	if err != nil {
		log.Println(command, err.Error())
	}
	result <- resultStruct
}

var ()

type Inventory struct {
	Urls   interface{}
	Flags  interface{}
	Values interface{}
}

func GenerateInventory(m interface{}) string {
	f, err := ioutil.TempFile("/tmp", "inventory-") // in Go version older than 1.17 you can use ioutil.TempFile
	if err != nil {
		log.Fatal(err)
	}

	// close and remove the temporary file at the end of the program
	defer f.Close()

	// write data to the temporary file
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	if _, err := f.Write(b); err != nil {
		log.Fatal(err)
	}
	return f.Name()
}

func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.Reply, error) {
	reply := pb.Reply{}
	result := make(chan []*pb.Result)

	filename := GenerateInventory(in.Data)
	defer os.Remove(filename)

	command := fmt.Sprintf("python3 %s/%s.py ping %s", *directory, in.Service, filename)
	go RunCheckers(command, result)

	reply.Results = <-result
	reply.Service = in.Service

	return &reply, nil
}

func (s *server) Put(ctx context.Context, in *pb.PutRequest) (*pb.Reply, error) {
	result := make(chan []*pb.Result)
	reply := pb.Reply{}

	filename := GenerateInventory(in.Data)
	defer os.Remove(filename)

	command := fmt.Sprintf("python3 %s/%s.py put %s %s", *directory, in.Service, in.Name, filename)
	go RunCheckers(command, result)

	reply.Results = <-result
	reply.Service = in.Service

	return &reply, nil
}

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.Reply, error) {
	result := make(chan []*pb.Result)
	reply := pb.Reply{}

	filename := GenerateInventory(in.Data)
	defer os.Remove(filename)

	command := fmt.Sprintf("python3 %s/%s.py get %s %s", *directory, in.Service, in.Name, filename)
	go RunCheckers(command, result)

	reply.Results = <-result
	reply.Service = in.Service

	return &reply, nil
}

func (s *server) Exploit(ctx context.Context, in *pb.ExploitRequest) (*pb.Reply, error) {
	result := make(chan []*pb.Result)
	reply := pb.Reply{}

	filename := GenerateInventory(in.Data)
	defer os.Remove(filename)

	command := fmt.Sprintf("python3 %s/%s.py exploit %s %s", *directory, in.Service, in.Name, filename)
	go RunCheckers(command, result)

	reply.Results = <-result
	reply.Service = in.Service

	return &reply, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCheckerServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}