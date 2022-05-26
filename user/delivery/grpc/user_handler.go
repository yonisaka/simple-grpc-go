package grpc

import (
	"context"
	"errors"
	"log"
	server "simple-grpc-go/config"
	pb "simple-grpc-go/user/delivery/grpc/protos"

	"github.com/jinzhu/gorm"
)

var users []*pb.User

type UserDataServer struct {
	pb.UnimplementedUserDataServer
}

func (s *UserDataServer) GetUsers(in *pb.Empty, stream pb.UserData_GetUsersServer) error {
	log.Printf("Received: %v", in)
	if err := server.DB.Find(&users).Error; err != nil {
		log.Fatalln(err)
	}
	for _, user := range users {
		if err := stream.Send(user); err != nil {
			return err
		}
	}
	return nil
}

func (s *UserDataServer) GetUser(ctx context.Context, in *pb.Id) (*pb.User, error) {
	log.Printf("Received: %v", in)
	res := &pb.User{}
	if err := server.DB.Find(&res).Where("id = ?", in).Error; err != nil {
		log.Fatalln(err)
	}
	return res, nil
}

func (s *UserDataServer) CreateUser(ctx context.Context, in *pb.User) (*pb.User, error) {
	log.Printf("Received: %v", in)
	if err := server.DB.Debug().Model(&pb.User{}).Create(&in).Error; err != nil {
		return nil, err
	}
	return in, nil
}

func (s *UserDataServer) UpdateUser(ctx context.Context, in *pb.User) (*pb.Status, error) {
	log.Printf("Received: %v", in)
	res := pb.Status{}
	if err := server.DB.Debug().Model(&pb.User{}).Where("id = ?", in.Id).Updates(pb.User{Name: in.Name, Email: in.Email, Age: in.Age}).Error; err != nil {
		return nil, err
	}
	res.Value = 1
	return &res, nil
}

func (s *UserDataServer) DeleteUser(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	log.Printf("Received: %v", in)
	res := pb.Status{}
	err := server.DB.Debug().Model(&pb.User{}).Where("id = ?", in.GetValue()).Take(&pb.User{}).Delete(&pb.User{})
	if err.Error != nil {
		res.Value = 0
		if gorm.IsRecordNotFoundError(server.DB.Error) {
			return &res, errors.New("User not found")
		}
		return &res, server.DB.Error
	}
	res.Value = 1
	return &res, nil
}