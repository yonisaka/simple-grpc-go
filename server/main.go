package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"

	pb "simple-grpc-go/user/delivery/grpc/protos"

	userDeliveryGrpc "simple-grpc-go/user/delivery/grpc"
	userRepo "simple-grpc-go/user/repository"
	userUcase "simple-grpc-go/user/usecase"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

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
	dbDriver := viper.GetString(`database.driver`)
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(dbDriver, dsn)
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
	}
	defer dbConn.Close()
	
	ur := userRepo.NewMysqlUserRepository(dbConn)
	uc := userUcase.NewUserUsecase(ur)
	list, err := net.Listen("tcp", viper.GetString("server.address"))
	if err != nil {
		fmt.Println("SOMETHING HAPPEN")
	}

	server := grpc.NewServer()
	userDeliveryGrpc.NewUserServerGrpc(server, uc)
	fmt.Println("Server Run at ", viper.GetString("server.address"))

	err = server.Serve(list)
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}
}