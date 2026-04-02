package tools

import "mcp-itinerary/mcp"

func HandleMCPRequest(req mcp.Request) mcp.Response {
	switch req.Method {

	case "initialize":
		return success(req.ID, map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{}, // 👈 must be object, not bool
			},
			"serverInfo": map[string]interface{}{
				"name":    "itinerary-mcp",
				"version": "1.0.0",
			},
		})

	case "tools/list":
		return success(req.ID, map[string]interface{}{
			"tools": ListTools(),
		})

	case "tools/call":
		return HandleToolCall(req)

	default:
		return success(req.ID, nil)
	}
}
