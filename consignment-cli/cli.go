package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/bobcats/shipper/consignment-service/proto/consignment"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
)

const (
	address              = "localhost:50051"
	defaultFilename      = "consignment.json"
	defaultTokenFilename = "token"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func readTokenFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return "", err
	}

	return string(data), err
}

func main() {
	cmd.Init()
	// Set up a connection to the server.
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)

	// Contact the server and print out its response.
	file := defaultFilename
	tokenFile := defaultTokenFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	if len(os.Args) > 2 {
		tokenFile = os.Args[2]
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	token, err := readTokenFile(tokenFile)
	if err != nil {
		log.Fatalf("Could not read token file: %v", err)
	}

	ctx := metadata.NewContext(context.Background(), map[string]string{
		"token": token,
	})

	r, err := client.CreateConsignment(ctx, consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)
	getAll, err := client.GetConsignments(ctx, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignments: %v", err)
	}

	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
