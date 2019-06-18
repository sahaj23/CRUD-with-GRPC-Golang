package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"google.golang.org/grpc"

	"tasks/taskuser/prototype"
)

func TestCreateUserProfile(t *testing.T) {
	conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()
	c := prototype.NewUserProfilesClient(conn)
	req := &prototype.CreateUserProfileRequest{
		UserProfile: &prototype.UserProfile{
			Id:        "6",
			FirstName: "Sahaj",
			LastName:  "Khandelwal",
			Email:     "sahaj@gmail.com",
		},
	}
	res, err := c.CreateUserProfile(context.Background(), req)
	fmt.Println(res)
}
