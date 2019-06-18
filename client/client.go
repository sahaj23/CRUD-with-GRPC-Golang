package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"tasks/taskuser/prototype"
)

func main() {
	fmt.Println("Helllo i am  a client")
	conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()
	c := prototype.NewUserProfilesClient(conn)

	//CREATE USER PROFILE

	req := &prototype.CreateUserProfileRequest{
		UserProfile: &prototype.UserProfile{
			Id:        "",
			FirstName: "Sahaj",
			LastName:  "Khandelwal",
			Email:     "sahaj@gmail.com",
		},
	}
	res, err := c.CreateUserProfile(context.Background(), req)

	//GET USER PROFILE

	// req := &prototype.GetUserProfileRequest{
	// 	Id: "6",
	// }
	// res, err := c.GetUserProfile(context.Background(), req)

	//UPDATE USER PROFILES

	// req := &prototype.UpdateUserProfileRequest{
	// 	UserProfile: &prototype.UserProfile{
	// 		Id:        "31",
	// 		FirstName: "Sahaj",
	// 		LastName:  "Khandelwal",
	// 		Email:     "sahaj@gmail.com",
	// 	},
	// }
	// res, err := c.UpdateUserProfile(context.Background(), req)

	//LIST USER PROFILES

	//req := &prototype.ListUsersProfilesRequest{Query: "a"}
	//res, err := c.ListUsersProfiles(context.Background(), req)

	//DELETE USER PROFILE

	//req := &prototype.DeleteUserProfileRequest{Id: "31"}
	// res, err := c.DeleteUserProfile(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while creating user profile %v", err)
	}
	fmt.Println(res)
}
