package utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapIDsToObjectIDs(ids []string) ([]primitive.ObjectID, error) {
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid object id: %s", id)
		}
		objectIDs = append(objectIDs, objectID)
	}

	return objectIDs, nil
}
