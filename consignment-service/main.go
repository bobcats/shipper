package main

import (
	"context"
	"log"
	"net"

	pb "github.com/bobcats/shipper/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// ConsignmentRepository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

// Create - Creates a consignment in the Consignment repository
func (repo *ConsignmentRepository) Create(c *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, c)
	repo.consignments = updated
	return c, nil
}

func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo Repository
}

// CreateConsignment - we created one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(req)

	if err != nil {
		return nil, err
	}

	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()
	return &pb.Response{Consignments: consignments}, nil
}

func main() {
	repo := &ConsignmentRepository{}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterShippingServiceServer(s, &service{repo})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
