package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/diorgenesbk/go-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRCPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	AddUserVerbose(client)
}

func AddUser(client pb.UserServiceClient) {
	request := &pb.User{
		Id:    "0",
		Name:  "Ariane",
		Email: "ariane.couto@gmail.com",
	}

	response, err := client.AddUser(context.Background(), request)

	if err != nil {
		log.Fatalf("Could not make gRCPC request: %v", err)
	}

	fmt.Println(response)
}

func AddUserVerbose(client pb.UserServiceClient) {
	request := &pb.User{
		Id:    "0",
		Name:  "Ariane",
		Email: "ariane.couto@gmail.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), request)

	if err != nil {
		log.Fatalf("Could not make gRCPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}

		fmt.Println("Status: ", stream.Status, " - ", stream.User)

	}
}
