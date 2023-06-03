package mongodb

import (
	"context"
	"fmt"

	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/globals"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteAuthByID(ctx context.Context, id string) (bool, error) {
	// Convert the ID string to an ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("invalid ID format: %v", err)
	}

	// Create a filter for the ID field
	filter := bson.M{"_id": oid}

	// Delete the document matching the filter
	result, err := globals.AuthDb.DeleteOne(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to delete auth: %v", err)
	}

	// Check if the document exists
	if result.DeletedCount == 0 {
		return false, fmt.Errorf("auth with ID %s not found", id)
	}

	return true, nil
}
