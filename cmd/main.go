package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/config"
	mongodb "github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/db"
	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/globals"
	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/messaging"
	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/pb"
	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/services"
	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/utils"
	"google.golang.org/grpc"
)

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
	globals.Db = mongodb.ConnectToMongoDB(mongodbURL)

	// Bind our collection to our global variable for use in other methods
	globals.AuthDb = globals.Db.Database(c.MongoDBDb).Collection(c.MongoDBCollection)

	// Construct the RabbitMQ URL
	globals.RabbitMQUrl = fmt.Sprintf("amqps://%s:%s@rattlesnake.rmq.cloudamqp.com/%s", c.RabbitMQUser, c.RabbitMQPwd, c.RabbitMQUser)

	//Connect to RabbitMQ
	fmt.Println("Connecting to RabbitMQ...")
	conn, err := messaging.ConnectToRabbitMQ(globals.RabbitMQUrl)
	if err != nil {
		log.Fatalf("Can't connect to RabbitMQ: %s", err)
	}

	// Start listening for messages RabbitMQ
	go messaging.ConsumeMessage(conn, "auth_queue", messaging.HandleMessage)

	s := services.Server{
		DB:  globals.AuthDb,
		Jwt: jwt,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
