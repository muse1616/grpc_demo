package main

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"grpc_server/pb"
	"time"
)

var Employees = []pb.Employee{
	{
		Id:           1,
		No:           1994,
		FirstName:    "Chandler",
		LastName:     "Bing",
		MonthSalary:  &pb.MonthSalary{
			Basic: 5000,
			Bonus: 125.5,
		},
		Status:       pb.EmployeeStatus_NORMAL,
		LastModified: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
	{
		Id:           2,
		No:           1995,
		FirstName:    "Leo",
		LastName:     "Messi",
		MonthSalary:  &pb.MonthSalary{
			Basic: 4000,
			Bonus: 225.5,
		},
		Status:       pb.EmployeeStatus_RESIGNED,
		LastModified: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
}
