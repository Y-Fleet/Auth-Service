package main

import (
	model "auth_service/Models"
	service "auth_service/Service"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/Yfleet/shared_proto/api"
)

type server struct {
	pb.UnimplementedAuthServiceServer
}
type server2 struct {
}

func (s *server) CheckUsers(ctx context.Context, req *pb.CheckUsersRequest) (*pb.CheckUsersResponse, error) {
	fmt.Println(req)

	dbAuth, cancel := ConToDb()
	defer cancel()

	collection := dbAuth.Database("Users").Collection("Users")
	var result model.User
	filter := bson.M{"Login": req.Login, "Password": req.Password}

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.CheckUsersResponse{Message: "Denied"}, nil
		} else {
			log.Printf("Error processing login request: %v", err)
			return nil, status.Errorf(codes.Internal, "Error with this service")
		}
	}
	fmt.Println("given you're cookie")
	return &pb.CheckUsersResponse{
		Message:      "Welcome",
		Token:        service.GenJwtToken(req.Login),
		TokenRefresh: service.GenRefreshToken(req.Login),
	}, nil
}

func ConToDb() (*mongo.Client, context.CancelFunc) {
	const connectionString = "mongodb://Auth:Auth@localhost:27020"
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	fmt.Println("Connected successfully to MongoDB")
	return client, cancel
}

func (s *server) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	fmt.Println(req)
	return &pb.RefreshTokenResponse{Token: service.GeneTokenFromRefreshToken(req.RefreshToken)}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})

	reflection.Register(s)
	log.Println("Starting microservice on :50055")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
