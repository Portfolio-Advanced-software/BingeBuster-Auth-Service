package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	mongodb "github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/db"
	"github.com/Portfolio-Advanced-software/BingeBuster-Auth-Service/pkg/globals"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Message struct {
	UserId string `json:"user_id"`
	Action string `json:"action"`
}

func HandleMessage(body []byte) error {
	jsonStr := string(body)
	var msg Message
	err := json.Unmarshal([]byte(jsonStr), &msg)
	if err != nil {
		log.Println("Failed to unmarshal JSON:", err)
		return err
	}

	switch msg.Action {
	case "deleteAllRecords":
		_, err := mongodb.DeleteAuthByID(context.Background(), msg.UserId)
		if err != nil {
			log.Println("Failed to delete all records:", err)
		}
	case "saveRecord":
		// Insert the data into the database, result contains the newly generated Object ID for the new document
		_, err := globals.AuthDb.InsertOne(globals.MongoCtx, msg)
		// check for potential errors
		if err != nil {
			// return internal gRPC error to be handled later
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Internal error: %v", err),
			)
		}
	default:
		fmt.Println("Unknown action:", msg.Action)
	}

	return nil
}
