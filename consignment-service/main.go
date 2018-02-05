package main

import (
	"log"

	pb "github.com/bobcats/shipper/consignment-service/proto/consignment"
	vesselProto "github.com/bobcats/shipper/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
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
	repo         Repository
	vesselClient vesselProto.VesselServiceClient
}

// CreateConsignment - we created one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("Found vessel %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}
	req.VesselId = vesselResponse.Vessel.Id

	consignment, err := s.repo.Create(req)

	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &ConsignmentRepository{}

	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
