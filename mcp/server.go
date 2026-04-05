package mcp

import (
	"bufio"
	"encoding/json"
	"os"
)

type Handler func(Request) Response

type Server struct {
	Handler Handler
}

func (s *Server) Start() {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			// Return a JSON-RPC parse error to the client.
			writeResponse(Response{
				JSONRPC: "2.0",
				ID:      nil,
				Error: map[string]interface{}{
					"code":    -32700,
					"message": "Parse error",
				},
			})
			continue
		}

		resp := s.Handler(req)

		// Notifications (e.g. notifications/initialized) return an empty Response.
		// Do NOT write anything back — notifications have no response.
		if resp.JSONRPC == "" {
			continue
		}

		writeResponse(resp)
	}
}

// writeResponse serialises resp as JSON and writes it to stdout followed by a newline.
// Using os.Stdout.Write instead of fmt.Println avoids double-newline on some platforms.
func writeResponse(resp Response) {
	bytes, _ := json.Marshal(resp)
	bytes = append(bytes, '\n')
	os.Stdout.Write(bytes) //nolint:errcheck
}
