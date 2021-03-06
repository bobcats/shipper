package main

import (
	"errors"

	pb "github.com/bobcats/shipper/user-service/proto/user"
	micro "github.com/micro/go-micro"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

const topic = "user.created"

type service struct {
	repo         Repository
	tokenService Authable
	Publisher    micro.Publisher
}

func (s *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := s.repo.Get(req.Id)
	if err != nil {
		return err
	}

	res.User = user
	return nil
}

func (s *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}

	res.Users = users
	return nil
}

func (s *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := s.tokenService.Encode(user)

	if err != nil {
		return err
	}

	res.Token = token
	return nil
}

func (s *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPass)
	if err := s.repo.Create(req); err != nil {
		return err
	}

	res.User = req

	err = s.Publisher.Publish(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {
	claims, err := s.tokenService.Decode(req.Token)

	if err != nil {
		return err
	}

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	res.Valid = true
	return nil
}
