package tools

import (
	"mcp-itinerary/mcp"
	"mcp-itinerary/service"
)

func ListTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "create_itinerary",
			"description": "Create a travel itinerary",
			"inputSchema": map[string]interface{}{ // ✅ FIXED
				"type": "object",
				"properties": map[string]interface{}{
					"user_id":     map[string]string{"type": "string"},
					"destination": map[string]string{"type": "string"},
					"days":        map[string]string{"type": "number"},
					"budget":      map[string]string{"type": "number"},
				},
				"required": []string{"user_id", "destination"}, // ✅ good practice
			},
		},
		{
			"name":        "get_itinerary",
			"description": "Get itinerary by ID",
			"inputSchema": map[string]interface{}{ // ✅ FIXED
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]string{"type": "string"},
				},
				"required": []string{"id"}, // ✅ good practice
			},
		},
	}
}

func HandleToolCall(req mcp.Request) mcp.Response {
	name := req.Params["name"].(string)
	args := req.Params["arguments"].(map[string]interface{})

	switch name {

	case "create_itinerary":
		result := service.CreateItinerary(args)
		return success(req.ID, result)

	case "get_itinerary":
		result := service.GetItinerary(args)
		return success(req.ID, result)

	default:
		return success(req.ID, nil)
	}
}

func success(id interface{}, result interface{}) mcp.Response {
	return mcp.Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}
