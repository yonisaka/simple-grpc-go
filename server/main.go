package main

import (
	"log"
	"net"

	server "simple-grpc-go/config"
	_userGrpcDelivery "simple-grpc-go/user/delivery/grpc"

	pb "simple-grpc-go/user/delivery/grpc/protos"
	// _userRepo "simple-grpc-go/user/repository"
	// _userUcase "simple-grpc-go/user/usecase"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const port = ":50051"

var users []*pb.User

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
	server.InitializeConnDB()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// ur := _userRepo.NewUserRepository(server.DB)
	// timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	// uc := _userUcase.NewUserUsecase(ur, timeoutContext)
	// _userGrpcDelivery.NewUserHandler(uc)

	pb.RegisterUserDataServer(s, &_userGrpcDelivery.UserDataServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}