package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"grpc_client/pb"
	"io"
	"log"
	"os"
	"time"
)

const port = ":5001"

func main() {
	// 记得设置goland环境变量 GODEBUG = x509ignoreCN=0
	creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		log.Fatalln(err.Error())
	}
	options := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost"+port, options...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 关闭连接
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}(conn)

	client := pb.NewEmployeeServiceClient(conn)
	//GetByNo(client)
	//GetAll(client)
	//AddPhoto(client)
	SaveAll(client)
}
func SaveAll(client pb.EmployeeServiceClient) {
	var employees = []pb.Employee{
		{
			Id:        124,
			No:        1994,
			FirstName: "Chandler",
			LastName:  "Bing",
			MonthSalary: &pb.MonthSalary{
				Basic: 5000,
				Bonus: 125.5,
			},
			Status: pb.EmployeeStatus_NORMAL,
			LastModified: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
		{
			Id:        123,
			No:        1995,
			FirstName: "Leo",
			LastName:  "Messi",
			MonthSalary: &pb.MonthSalary{
				Basic: 4000,
				Bonus: 225.5,
			},
			Status: pb.EmployeeStatus_RESIGNED,
			LastModified: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	}

	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}
	finishChannel := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				finishChannel <- struct{}{}
				break
			}
			if err != nil {
				log.Fatalln(err.Error())
			}
			fmt.Println(res.Employee)
		}
	}()
	for _, e := range employees {
		err := stream.Send(&pb.EmployeeRequest{Employee: &e})
		if err!=nil{
			log.Fatalln(err.Error())
		}
	}
	err = stream.CloseSend()
	if err != nil {
		log.Fatalln(err.Error())
	}
	<-finishChannel
}
func AddPhoto(client pb.EmployeeServiceClient) {
	imgFile, err := os.Open("1.jpg")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer imgFile.Close()
	md := metadata.New(map[string]string{"no": "1994"})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.AddPhoto(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	for {
		chunk := make([]byte, 128*1024)
		chunkSize, err := imgFile.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		// 最后一块
		if chunkSize < len(chunk) {
			chunk = chunk[:chunkSize]
		}
		stream.Send(&pb.AddPhotoRequest{Data: chunk})

	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(res.IsOk)

}
func GetByNo(client pb.EmployeeServiceClient) {
	res, err := client.GetByNo(context.Background(), &pb.GetByNoRequest{No: 1994})
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(res.Employee)
}
func GetAll(client pb.EmployeeServiceClient) {
	stream, err := client.GetAll(context.Background(), &pb.GetAllRequest{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF { // 读取完毕
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println(res.Employee)
	}
}
