git@<ssh-alias>:<github-username>/<repo>.git
git remote set-url origin git@github-thecodelogs:thecodelogs/bonaire-mcp-server.git

# MCP Itinerary Server (Go) — Detailed Documentation

## 📌 Overview

This project implements a **minimal Model Context Protocol (MCP) server in Go** for an itinerary system.

The goal is to:

- Expose backend functionality as **AI-callable tools**
- Follow proper **JSON-RPC over stdio**
- Maintain **clean architecture (no circular dependencies)**

---

## 🧠 High-Level Architecture

```
          ┌───────────────┐
          │   AI Client   │ (Inspector / LLM)
          └──────┬────────┘
                 │ JSON-RPC (stdio)
                 ▼
        ┌────────────────────┐
        │     MCP Layer      │ (protocol handling)
        └────────┬───────────┘
                 ▼
        ┌────────────────────┐
        │    Tools Layer     │ (MCP adapters)
        └────────┬───────────┘
                 ▼
        ┌────────────────────┐
        │   Service Layer    │ (business logic)
        └────────────────────┘
```

---

## 📂 Project Structure

```
mcp-itinerary/
├── main.go
├── mcp/
│   ├── server.go
│   └── types.go
├── tools/
│   ├── router.go
│   └── itinerary.go
├── service/
│   └── itinerary_service.go
```

---

# 🔍 File-by-File Breakdown

---

## 1. `main.go` — Entry Point

### 🎯 Purpose

- Bootstraps the MCP server
- Injects the request handler (dependency injection)

### 🧩 Code Responsibility

```go
server := mcp.Server{
    Handler: tools.HandleMCPRequest,
}
```

### 💡 Key Concept

- **Dependency Injection**
- `mcp` package does NOT know about tools
- Avoids circular dependency

---

## 2. `mcp/types.go` — Protocol Models

### 🎯 Purpose

Defines JSON-RPC structures.

### 🧩 Structures

#### Request

```go
type Request struct {
    ID     interface{}
    Method string
    Params map[string]interface{}
}
```

#### Response

```go
type Response struct {
    JSONRPC string
    ID      interface{}
    Result  interface{}
    Error   interface{}
}
```

### 💡 Key Concept

- Mirrors JSON-RPC 2.0 spec
- Keeps protocol logic isolated

---

## 3. `mcp/server.go` — Core MCP Engine

### 🎯 Purpose

- Handles stdio communication
- Runs infinite loop
- Converts raw input → structured request → response

### 🔁 Flow

```
stdin → read → parse JSON → call handler → write stdout
```

### 🧩 Key Function

```go
func (s *Server) Start()
```

### 🔍 Internals

#### Step 1: Read input

```go
line, err := reader.ReadBytes('\n')
```

#### Step 2: Parse JSON

```go
json.Unmarshal(line, &req)
```

#### Step 3: Call handler

```go
resp := s.Handler(req)
```

#### Step 4: Return response

```go
fmt.Println(string(bytes))
```

### ⚠️ Important Rules

- NEVER log to stdout ❌
- Use stderr for debugging ✅
- Must run continuously (no exit)

---

## 4. `tools/router.go` — MCP Method Router

### 🎯 Purpose

Routes incoming MCP methods:

- `initialize`
- `tools/list`
- `tools/call`

---

### 🔁 Flow

```
Request.Method → switch → appropriate handler
```

---

### 🧩 Handled Methods

---

### ✅ `initialize`

```go
case "initialize":
```

#### Purpose:

- First handshake with client

#### Returns:

```json
{
  "protocolVersion": "2024-11-05",
  "capabilities": {
    "tools": {}
  },
  "serverInfo": {
    "name": "itinerary-mcp",
    "version": "1.0.0"
  }
}
```

---

### ✅ `tools/list`

```go
case "tools/list":
```

#### Purpose:

- Tell AI what tools are available

#### Calls:

```go
ListTools()
```

---

### ✅ `tools/call`

```go
case "tools/call":
```

#### Purpose:

- Execute a tool

#### Calls:

```go
HandleToolCall(req)
```

---

## 5. `tools/itinerary.go` — Tool Definitions

### 🎯 Purpose

- Define tools (metadata)
- Execute tool logic

---

### 🧩 Function 1: `ListTools()`

Returns:

```json
[
  {
    "name": "create_itinerary",
    "description": "...",
    "input_schema": { ... }
  }
]
```

---

### 🧩 Function 2: `HandleToolCall()`

```go
name := req.Params["name"]
args := req.Params["arguments"]
```

---

### 🔁 Flow

```
tool name → switch → call service layer
```

---

### Supported Tools

---

### ✅ `create_itinerary`

Calls:

```go
service.CreateItinerary(args)
```

---

### ✅ `get_itinerary`

Calls:

```go
service.GetItinerary(args)
```

---

## 6. `service/itinerary_service.go` — Business Logic

### 🎯 Purpose

Core logic of your application

---

### 🧩 Storage (Temporary)

```go
var store = make(map[string]interface{})
```

👉 This is **in-memory DB (for now)**

---

### 🧩 Function: `CreateItinerary`

#### Steps:

1. Generate ID
2. Build itinerary object
3. Store in map
4. Return result

---

### 🧩 Function: `GetItinerary`

#### Steps:

1. Read ID from input
2. Lookup in map
3. Return result / error

---

# 🔁 End-to-End Flow

---

## 🧪 Example: Create Itinerary

### Step 1: Client Request

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_itinerary",
    "arguments": {
      "user_id": "123",
      "destination": "Goa",
      "days": 3,
      "budget": 20000
    }
  }
}
```

---

### Step 2: Flow

```
MCP Server
 → Router
   → tools.HandleToolCall
     → service.CreateItinerary
       → store data
 → Response returned
```

---

### Step 3: Response

```json
{
  "id": "...",
  "result": {
    "id": "uuid",
    "destination": "Goa",
    ...
  }
}
```

---

# ⚠️ Design Decisions (Important)

---

## ❌ Why NOT import tools inside MCP?

Because:

```
mcp → tools
tools → mcp
```

👉 Causes circular dependency

---

## ✅ Solution Used

- `mcp` defines **Handler interface**
- `tools` implements it
- `main` wires them together

👉 Clean architecture

---

# 🧠 Key Concepts You Practiced

---

## 1. JSON-RPC over stdio

- Core of MCP communication

---

## 2. Tool-based architecture

- Backend functions → AI tools

---

## 3. Dependency Injection

- Decouples layers

---

## 4. Layered Design

- Protocol
- Adapter
- Business logic

---

## 5. Schema-driven API

- MCP is strict (unlike REST)

---

# 🚀 Current Limitations

- ❌ No database (in-memory only)
- ❌ No validation
- ❌ No error handling (structured)
- ❌ No authentication
- ❌ No logging system

---

# 🔥 Next Step (Step 2)

We will upgrade to:

- ✅ PostgreSQL integration
- ✅ Input validation
- ✅ Proper error responses
- ✅ Clean DTOs

---

# 💬 Summary

You have built:

✔ A working MCP server
✔ Tool-based architecture
✔ AI-callable backend
✔ Clean, scalable structure

---

👉 This is **foundation of AI backend systems**

---

When ready, say:

**"step 2"** 🚀
