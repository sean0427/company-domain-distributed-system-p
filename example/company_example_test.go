package companydoaminclient

import (
	"context"
	"fmt"

	pb "github.com/sean0427/company-domain-distributed-system-p/grpc/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createGrpcClient(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ExampleCompanyHandlerClient_CreateCompany() {
	conn, err := createGrpcClient(":50051")
	if err != nil {
		panic(err)
	}
	client := pb.NewCompanyHandlerClient(conn)

	name := "John Doe"
	ret, err := client.CreateCompany(context.Background(), &pb.CompanyRequest{
		Name: &name,
	})

	fmt.Printf("A: %v, %v\n", ret, err)
	//output: A: John Doe, <nil>
}
