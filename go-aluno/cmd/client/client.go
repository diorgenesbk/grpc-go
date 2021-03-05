package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUserStreamBoth(client)
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

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "Didico",
			Email: "didico@teste.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "Ariane",
			Email: "ariane@teste.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "Preta",
			Email: "preta@teste.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "Pipoca",
			Email: "pipoca@teste.com",
		},
		&pb.User{
			Id:    "5",
			Name:  "Bolinho",
			Email: "bolinho@teste.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "Didico",
			Email: "didico@teste.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "Ariane",
			Email: "ariane@teste.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "Preta",
			Email: "preta@teste.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "Pipoca",
			Email: "pipoca@teste.com",
		},
		&pb.User{
			Id:    "5",
			Name:  "Bolinho",
			Email: "bolinho@teste.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user", req.GetName())
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v ", err)
				break
			}

			fmt.Println("Recebendo user ", res.GetUser().GetName(), " com status: ", res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
