package main

import (
	"mcp-itinerary/mcp"
	"mcp-itinerary/tools"
)

func main() {
	server := mcp.Server{
		Handler: tools.HandleMCPRequest, // 👈 inject here
	}

	server.Start()
}
