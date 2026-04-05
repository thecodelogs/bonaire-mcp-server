package tools

import (
	"encoding/json"
	"fmt"
	"mcp-itinerary/mcp"
	"mcp-itinerary/service"
)

func ListTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "create_itinerary",
			"description": "Create a travel itinerary for a user given a destination, number of days, and optional budget.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"user_id":     map[string]string{"type": "string", "description": "Unique identifier of the user"},
					"destination": map[string]string{"type": "string", "description": "Travel destination (e.g. Goa, Paris)"},
					"days":        map[string]string{"type": "number", "description": "Number of days for the trip"},
					"budget":      map[string]string{"type": "number", "description": "Total budget in local currency"},
				},
				"required": []string{"user_id", "destination"},
			},
		},
		{
			"name":        "get_itinerary",
			"description": "Retrieve a previously created itinerary by its unique ID.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]string{"type": "string", "description": "Unique ID of the itinerary to retrieve"},
				},
				"required": []string{"id"},
			},
		},
	}
}

func HandleToolCall(req mcp.Request) mcp.Response {
	name, _ := req.Params["name"].(string)
	args, _ := req.Params["arguments"].(map[string]interface{})

	switch name {

	case "create_itinerary":
		result := service.CreateItinerary(args)
		return toolSuccess(req.ID, result)

	case "get_itinerary":
		result := service.GetItinerary(args)
		return toolSuccess(req.ID, result)

	default:
		return toolError(req.ID, fmt.Sprintf("Unknown tool: %s", name))
	}
}

// toolSuccess wraps a result in the MCP-required content array format.
func toolSuccess(id interface{}, result interface{}) mcp.Response {
	text, _ := json.Marshal(result)
	return mcp.Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": string(text)},
			},
			"isError": false,
		},
	}
}

// toolError returns a tool-level error (NOT a JSON-RPC protocol error).
// Per spec: tool execution errors use isError:true inside the result, not the error field.
func toolError(id interface{}, message string) mcp.Response {
	return mcp.Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]string{
				{"type": "text", "text": message},
			},
			"isError": true,
		},
	}
}
