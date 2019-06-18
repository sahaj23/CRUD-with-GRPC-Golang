package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	empty "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/satori/uuid"
	"google.golang.org/grpc"

	"tasks/taskuser/prototype"
	pb "tasks/taskuser/prototype"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "sahaj231197"
	dbname   = "postgres"
)

type server struct {
	conn *sql.DB
}

func (connection *server) CreateUserProfile(ctx context.Context, req *pb.CreateUserProfileRequest) (*pb.UserProfile, error) {
	db := connection.conn
	id, err := uuid.NewV4()
	if err != nil {
		errors.Wrap(err, "Couldn't generate uuid")
	}
	req.UserProfile.Id = id.String()
	firstname := req.GetUserProfile().GetFirstName()
	lastname := req.GetUserProfile().GetLastName()
	email := req.GetUserProfile().GetEmail()
	sqlStatement := `INSERT INTO "user" ( first_name, last_name, email,id) VALUES ($1, $2, $3, $4)`
	if _, err := db.Exec(sqlStatement, firstname, lastname, email, id); err != nil {
		return nil, errors.Wrap(err, "User couldn't be inserted")
	}
	return req.UserProfile, nil
}
func (connection *server) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.UserProfile, error) {
	db := connection.conn
	id := req.GetId()
	sqlStatement := `select * from "user" where id=$1`
	var first, last, email, uid string
	err := db.QueryRow(sqlStatement, id).Scan(&first, &last, &email, &uid)
	if err != nil {
		errors.Wrap(err, "UserProfile couldn't be returned")
	}
	res := &pb.UserProfile{
		FirstName: first,
		LastName:  last,
		Email:     email,
		Id:        id,
	}
	return res, nil
}
func (connection *server) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UserProfile, error) {
	db := connection.conn
	sqlStatement := `UPDATE "user" SET first_name=$1, last_name=$2,email=$3 WHERE "id" =$4;`
	if _, err := db.Exec(sqlStatement, req.UserProfile.FirstName, req.UserProfile.LastName, req.UserProfile.Email, req.GetUserProfile().GetId()); err != nil {
		return nil, err
	}
	return req.UserProfile, nil
}
func (connection *server) ListUsersProfiles(ctx context.Context, req *pb.ListUsersProfilesRequest) (*pb.ListUsersProfilesResponse, error) {
	db := connection.conn
	id := req.GetQuery() + "%"
	sqlStatement := `select * from "user" where first_name like $1`
	result, err := db.Query(sqlStatement, id)
	defer result.Close()
	if err != nil {
		fmt.Println(err)
	}
	res := []*pb.UserProfile{}
	for result.Next() {
		var first, last, email, id string
		if err = result.Scan(&first, &last, &email, &id); err != nil {
			errors.Wrap(err, "Users couln't be listed")
		}
		u := pb.UserProfile{
			FirstName: first,
			LastName:  last,
			Email:     email,
			Id:        id,
		}
		res = append(res, &u)
	}
	ans := pb.ListUsersProfilesResponse{Profiles: res}
	return &ans, nil
}
func (connection *server) DeleteUserProfile(ctx context.Context, req *pb.DeleteUserProfileRequest) (*empty.Empty, error) {
	db := connection.conn
	id := req.GetId()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	sqlStatement := `delete from "user" where id=$1`
	if _, err = db.Exec(sqlStatement, id); err != nil {
		errors.Wrap(err, "User couldn't be deleted")
	}
	return &empty.Empty{}, nil
}
func main() {
	fmt.Println("Welcome to the server")
	lis, err := net.Listen("tcp", "0.0.0.0:50056")
	if err != nil {
		errors.Wrap(err, " Failed to listen the port")
	}
	s := grpc.NewServer()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		errors.Wrap(err, "Connection couldn't be opened")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		errors.Wrap(err, "Connection not established, ping didn't work")
	}
	prototype.RegisterUserProfilesServer(s, &server{db})
	if err := s.Serve(lis); err != nil {
		errors.Wrap(err, "Failed to server the listener")
	}
}
