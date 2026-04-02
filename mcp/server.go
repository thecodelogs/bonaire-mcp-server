package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
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
			continue
		}

		resp := s.Handler(req)
		bytes, _ := json.Marshal(resp)

		fmt.Println(string(bytes))
	}
}
