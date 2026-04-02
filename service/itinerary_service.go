package service

import (
	"time"

	"github.com/google/uuid"
)

var store = make(map[string]interface{})

func CreateItinerary(args map[string]interface{}) map[string]interface{} {
	id := uuid.New().String()

	itinerary := map[string]interface{}{
		"id":          id,
		"user_id":     args["user_id"],
		"destination": args["destination"],
		"days":        args["days"],
		"budget":      args["budget"],
		"activities":  []string{},
		"created_at":  time.Now().String(),
	}

	store[id] = itinerary
	return itinerary
}

func GetItinerary(args map[string]interface{}) interface{} {
	id := args["id"].(string)

	if val, ok := store[id]; ok {
		return val
	}

	return map[string]string{
		"error": "not found",
	}
}
