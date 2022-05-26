package domain

import (
	"context"
	pb "simple-grpc-go/user/delivery/grpc/protos"
)

// type User struct {
// 	ID        int64     `json:"id"`
// 	Name      string    `json:"name"`
// 	Email     string    `json:"email"`
// 	Age       int64     `json:"age"`
// }

type UserUsecase interface {
	Fetch(ctx context.Context, in *pb.Empty) ([]pb.User, error)
	GetByID(ctx context.Context, id int64) (pb.User, error)
	Update(ctx context.Context, user *pb.User) error
	Store(context.Context, *pb.User) error
	Delete(ctx context.Context, id int64) error
}

type UserRepository interface {
	Fetch(ctx context.Context, in *pb.Empty) (res []pb.User, err error)
	GetByID(ctx context.Context, id int64) (pb.User, error)
	Update(ctx context.Context, user *pb.User) error
	Store(ctx context.Context, user *pb.User) error
	Delete(ctx context.Context, id int64) error
}
