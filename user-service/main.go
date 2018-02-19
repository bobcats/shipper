package main

import (
	"fmt"
	"log"

	pb "github.com/bobcats/shipper/user-service/proto/user"
	"github.com/micro/go-micro"
)

func main() {
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	db.AutoMigrate(&pb.User{})
	repo := &UserRepository{db}

	tokenService := &TokenService{repo}

	srv := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	srv.Init()

	publisher := micro.NewPublisher("user.created", srv.Client())

	pb.RegisterUserServiceHandler(srv.Server(), &service{repo, tokenService, publisher})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
