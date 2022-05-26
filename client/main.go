package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "simple-grpc-go/user/delivery/grpc/protos"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := pb.NewUserDataClient(conn)

	// runGetUsers(client)
	// runGetUser(client, 2)
	// runCreateUser(client, "Momo", "momo@gmail.com", 12)
	// runUpdateUser(client, 10, "Momos", "momo@gmail.com", 12)
	runDeleteUser(client, 10)
}

func runGetUsers(client pb.UserDataClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Empty{}
	stream, err := client.GetUsers(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetUsers(_) = _, %v", client, err)
	}
	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetUsers(_) = _, %v", client, err)
		}
		log.Printf("Users: %v", row)
	}
}

func runGetUser(client pb.UserDataClient, userid int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: userid}
	res, err := client.GetUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetUser(_) = _, %v", client, err)
	}
	log.Printf("UserInfo: %v", res)
}

func runCreateUser(client pb.UserDataClient, name string, email string, age int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.User{Name: name, Email: email, Age: age}
	res, err := client.CreateUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.CreateUser(_) = _, %v", client, err)
	}
	if res != nil{
		log.Printf("CreateUser: %v", res)
	} else {
		log.Printf("CreateUser Failed")
	}
}

func runUpdateUser(client pb.UserDataClient, userid int64, name string, email string, age int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.User{Id: userid, Name: name, Email: email, Age: age}
	res, err := client.UpdateUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.UpdateUser(_) = _, %v", client, err)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("UpdateUser Success")
	} else {
		log.Printf("UpdateUser Failed")
	}
}

func runDeleteUser(client pb.UserDataClient, userid int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: userid}
	res, err := client.DeleteUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.DeleteUser(_) = _, %v", client, err)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("DeleteUser Success")
	} else {
		log.Printf("DeleteUser Failed")
	}
}
