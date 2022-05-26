package main

import (
	"context"
	"errors"
	"log"
	"net"

	. "simple-grpc-go/config"
	pb "simple-grpc-go/user/delivery/grpc/protos"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const port = ":50051"

var users []*pb.User

type userDataServer struct {
	pb.UnimplementedUserDataServer
}

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
	InitializeConnDB()
	// initUsers(DB)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterUserDataServer(s, &userDataServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initUsers(db *gorm.DB) {
	if err := db.Find(&users).Error; err != nil {
		log.Fatalln(err)
	}
}

func (s *userDataServer) GetUsers(in *pb.Empty, stream pb.UserData_GetUsersServer) error {
	log.Printf("Received: %v", in)
	if err := DB.Find(&users).Error; err != nil {
		log.Fatalln(err)
	}
	for _, user := range users {
		if err := stream.Send(user); err != nil {
			return err
		}
	}
	return nil
}

func (s *userDataServer) GetUser(ctx context.Context, in *pb.Id) (*pb.User, error) {
	log.Printf("Received: %v", in)
	res := &pb.User{}
	if err := DB.Find(&res).Where("id = ?", in).Error; err != nil {
		log.Fatalln(err)
	}
	return res, nil
}

func (s *userDataServer) CreateUser(ctx context.Context, in *pb.User) (*pb.User, error) {
	log.Printf("Received: %v", in)
	if err := DB.Debug().Model(&pb.User{}).Create(&in).Error; err != nil {
		return nil, err
	}
	return in, nil
}

func (s *userDataServer) UpdateUser(ctx context.Context, in *pb.User) (*pb.Status, error) {
	log.Printf("Received: %v", in)
	res := pb.Status{}
	if err := DB.Debug().Model(&pb.User{}).Where("id = ?", in.Id).Updates(pb.User{Name: in.Name, Email: in.Email, Age: in.Age}).Error; err != nil {
		return nil, err
	}
	res.Value = 1
	return &res, nil
}

func (s *userDataServer) DeleteUser(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	log.Printf("Received: %v", in)
	res := pb.Status{}
	err := DB.Debug().Model(&pb.User{}).Where("id = ?", in.GetValue()).Take(&pb.User{}).Delete(&pb.User{})
	if err.Error != nil {
		res.Value = 0
		if gorm.IsRecordNotFoundError(DB.Error) {
			return &res, errors.New("User not found")
		}
		return &res, DB.Error
	}
	res.Value = 1
	return &res, nil
}