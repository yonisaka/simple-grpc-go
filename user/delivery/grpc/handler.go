package grpc

import (
	"context"
	"log"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	models "simple-grpc-go/user"
	pb "simple-grpc-go/user/delivery/grpc/protos"
	_usecase "simple-grpc-go/user/usecase"

	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
)

func NewUserServerGrpc(gserver *grpc.Server, userUcase _usecase.UserUsecase) {

	userServer := &server{
		usecase: userUcase,
	}

	pb.RegisterUserHandlerServer(gserver, userServer)
	reflection.Register(gserver)
}

type server struct {
	usecase _usecase.UserUsecase
}

func (s *server) transformUserRPC(ur *models.User) *pb.User {

	if ur == nil {
		return nil
	}

	updated_at := &google_protobuf.Timestamp{
		Seconds: ur.UpdatedAt.Unix(),
	}
	created_at := &google_protobuf.Timestamp{
		Seconds: ur.CreatedAt.Unix(),
	}
	res := &pb.User{
		Id: ur.ID,
		Name: ur.Name,
		Email: ur.Email,
		Age: ur.Age,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}
	return res
}

func (s *server) transformUserData(ur *pb.User) *models.User {
	updated_at := time.Unix(ur.GetUpdatedAt().GetSeconds(), 0)
	created_at := time.Unix(ur.GetCreatedAt().GetSeconds(), 0)
	res := &models.User{
		ID: ur.Id,
		Name: ur.Name,
		Email: ur.Email,
		Age: ur.Age,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}
	return res
}

func (s *server) GetUsers(in *pb.FetchRequest, stream pb.UserHandler_GetUsersServer) error {
	cursor := ""
	num := int64(0)
	if in != nil {
		cursor = in.Cursor
		num = in.Num
	}
	list, _, err := s.usecase.Fetch(cursor, num)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, a := range list {
		ar := s.transformUserRPC(a)

		if err := stream.Send(ar); err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

func (s *server) GetUser(ctx context.Context, in *pb.SingleRequest) (*pb.User, error) {
	id := int64(0)
	if in != nil {
		id = in.Id
	}
	ar, err := s.usecase.GetByID(id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if ar == nil {
		return nil, models.NOT_FOUND_ERROR
	}

	res := s.transformUserRPC(ar)
	return res, nil
}

func (s *server) CreateUser(ctx context.Context, a *pb.User) (*pb.User, error) {
	ar := s.transformUserData(a)
	data, err := s.usecase.Store(ar)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	res := s.transformUserRPC(data)

	return res, nil
}

func (s *server) UpdateUser(c context.Context, ar *pb.User) (*pb.User, error) {
	a := s.transformUserData(ar)
	res, err := s.usecase.Update(a)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	l := s.transformUserRPC(res)
	return l, nil
}

func (s *server) DeleteUser(c context.Context, in *pb.SingleRequest) (*pb.DeleteResponse, error) {
	id := int64(0)
	if in != nil {
		id = in.Id
	}

	ok, err := s.usecase.Delete(id)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		return nil, err
	}
	resp := &pb.DeleteResponse{
		Status: "Not Oke To Delete",
	}
	if ok {
		resp.Status = "Succesfull To Delete"
	}

	return resp, nil
}