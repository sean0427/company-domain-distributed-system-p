package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	service "github.com/sean0427/company-domain-distributed-system-p"
	config "github.com/sean0427/company-domain-distributed-system-p/config"
	handler "github.com/sean0427/company-domain-distributed-system-p/grpc"
	pb "github.com/sean0427/company-domain-distributed-system-p/grpc/grpc"
	repository "github.com/sean0427/company-domain-distributed-system-p/postgressql"
)

func NewSQLDB() (*gorm.DB, error) {
	dsn, err := config.GetPostgresDNS()
	if err != nil {
		return nil, err
	}

	gormConfig := &gorm.Config{}
	conn, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

var (
	port = flag.Int("port", 50051, "The server port")
)

func startServer() {
	fmt.Println("Starting server...")

	db, err := NewSQLDB()
	if err != nil {
		panic(err)
	}

	r := repository.New(db)
	s := service.New(r)
	h := handler.New(s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
	netServer := grpc.NewServer()
	pb.RegisterCompanyHandlerServer(netServer, h)

	log.Printf("server listening at %v", lis.Addr())
	if err := netServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	fmt.Println("Stoping server...")
}

func main() {
	startServer()
}
