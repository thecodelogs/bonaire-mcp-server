package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const toursAPI = "https://stage-api.visit-bonaire.com/api/v2/gettours"

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

// FetchAllTours calls the Bonaire API and returns the list of available tours.
func FetchAllTours(_ map[string]interface{}) interface{} {
	resp, err := http.Get(toursAPI)
	if err != nil {
		return map[string]string{"error": fmt.Sprintf("failed to fetch tours: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return map[string]string{"error": fmt.Sprintf("API returned status %d", resp.StatusCode)}
	}

	var body struct {
		Tours []map[string]interface{} `json:"Tours"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return map[string]string{"error": fmt.Sprintf("failed to decode response: %v", err)}
	}

	return body.Tours
}
