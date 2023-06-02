package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/config"
	mongodb "github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/db"
	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/pb"
	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/services"
	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

var db *mongo.Client
var authdb *mongo.Collection
var mongoCtx context.Context

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "go-grpc-auth-svc",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	fmt.Println("Auth Svc on", c.Port)

	// Construct the RabbitMQ URL
	mongodbURL := fmt.Sprintf("mongodb+srv://%s:%s@%s", c.MongoDBUser, c.MongoDBPwd, c.MongoDBCluster)

	// Initialize MongoDb client
	fmt.Println("Connecting to MongoDB...")
	db = mongodb.ConnectToMongoDB(mongodbURL)

	// Bind our collection to our global variable for use in other methods
	authdb = db.Database(c.MongoDBDb).Collection(c.MongoDBCollection)

	s := services.Server{
		DB:  authdb,
		Jwt: jwt,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
