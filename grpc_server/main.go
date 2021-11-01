package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"grpc_server/pb"
	"io"
	"log"
	"net"
	"time"
)

const port = ":5001"

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err.Error())
	}
	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatalln(err.Error())
	}
	options := []grpc.ServerOption{grpc.Creds(creds)}
	server := grpc.NewServer(options...)
	pb.RegisterEmployeeServiceServer(server, new(employeeService))
	log.Println("gRPC Server Started..." + port)
	err = server.Serve(listen)
	if err != nil {
		return
	}
}

type employeeService struct {
}

func (e employeeService) GetByNo(ctx context.Context, request *pb.GetByNoRequest) (*pb.EmployeeResponse, error) {
	for _, e := range Employees {
		if request.No == e.No {
			return &pb.EmployeeResponse{
				Employee: &e,
			}, nil
		}
	}
	return nil, errors.New("Employee Not Found ")
}

func (e employeeService) GetAll(request *pb.GetAllRequest, server pb.EmployeeService_GetAllServer) error {
	for _, e := range Employees {
		err := server.Send(&pb.EmployeeResponse{
			Employee: &e,
		})
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (e employeeService) AddPhoto(server pb.EmployeeService_AddPhotoServer) error {
	metadata, ok := metadata.FromIncomingContext(server.Context())
	if ok {
		fmt.Printf("Employee:%s\n", metadata["no"][0])
	}
	var img []byte
	for {
		data, err := server.Recv()
		if err == io.EOF {
			fmt.Printf("File Size:%d\n", len(img))
			return server.SendAndClose(&pb.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		fmt.Printf("File received: %d\n", len(data.Data))
		img = append(img, data.Data...)
	}
}

func (e employeeService) Save(ctx context.Context, request *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	panic("implement me")
}

func (e employeeService) SaveAll(server pb.EmployeeService_SaveAllServer) error {
	for {
		empReq, err := server.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		Employees = append(Employees, *empReq.Employee)
		err = server.Send(&pb.EmployeeResponse{Employee: empReq.Employee})
		if err != nil {
			return err
		}
	}
	for _, emp := range Employees {
		fmt.Println(emp)
	}
	return nil
}
