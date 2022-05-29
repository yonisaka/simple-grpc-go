package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "simple-grpc-go/user/delivery/grpc/protos"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	address := viper.GetString(`server.address`)
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Unexpected Error", err)
	}

	defer conn.Close()
	c := pb.NewUserHandlerClient(conn)
	getUsers(c)
	// getUserById(c)
	// creteUser(c)
	// updateUser(c)
	// deleteUser(c)
}

func getUsers(c pb.UserHandlerClient) {
	f := &pb.FetchRequest{
		Num: 0,
		Cursor: "",
	}
	stream, err := c.GetUsers(context.Background(), f)
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Unexpected Error", err)
			break
		}
		fmt.Println("Data From Stream :  ", r.Name, r.Email)
	}
}

func getUserById(c pb.UserHandlerClient) {
	req := &pb.SingleRequest{
		Id: 2,
	}
	r, err := c.GetUser(context.Background(), req)
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}
	fmt.Println("Data From Single Request : ", r.Name, r.Email)
}

func creteUser(c pb.UserHandlerClient) {
	req := &pb.User{
		Name: "John",
		Email: "john@gmail.com",
		Age: 23,
	}
	r, err := c.CreateUser(context.Background(), req)
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}
	fmt.Println("User Created : ", r.Name)
}

func updateUser(c pb.UserHandlerClient) {
	req := &pb.User{
		Id : 11,
		Name: "Johns",
		Email: "john@gmail.com",
		Age: 23,
	}
	r, err := c.UpdateUser(context.Background(), req)
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}
	fmt.Println("User Updated : ", r.Name)
}

func deleteUser(c pb.UserHandlerClient) {
	req := &pb.SingleRequest{
		Id: 14,
	}
	r, err := c.DeleteUser(context.Background(), req)
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}
	fmt.Println("User Deleted : ", r.Status)
}