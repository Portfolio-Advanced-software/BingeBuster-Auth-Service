package globals

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Db             *mongo.Client
	AuthDb         *mongo.Collection
	MongoCtx       context.Context
	MongoDBUrl     string
	DbName         string
	CollectionName string
	RabbitMQUrl    string
)
