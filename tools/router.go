package tools

import "mcp-itinerary/mcp"

func HandleMCPRequest(req mcp.Request) mcp.Response {
	switch req.Method {

	case "initialize":
		return protocolSuccess(req.ID, map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "itinerary-mcp",
				"version": "1.0.0",
			},
		})

	case "notifications/initialized":
		// This is a notification (no id). Silently ignore — do not send a response.
		return mcp.Response{}

	case "tools/list":
		return protocolSuccess(req.ID, map[string]interface{}{
			"tools": ListTools(),
		})

	case "tools/call":
		return HandleToolCall(req)

	default:
		return protocolSuccess(req.ID, map[string]interface{}{})
	}
}

// protocolSuccess builds a standard JSON-RPC success response for protocol-level methods.
func protocolSuccess(id interface{}, result interface{}) mcp.Response {
	return mcp.Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}
