Deleted router.go
Edited main.go
Deleted types.go
Edited router.go
Edited itinerary.go
Edited README.md

Yes, these keys are required by the **MCP specification**. Here's what each one does:

```go
{
    "name":        "create_itinerary",       // ✅ REQUIRED — tool identifier
    "description": "Create a travel itinerary", // ✅ REQUIRED — LLM reads this to decide WHEN to call the tool
    "inputSchema": map[string]interface{}{   // ✅ REQUIRED — tells LLM what arguments to pass
        "type": "object",                    // ✅ REQUIRED — must always be "object"
        "properties": map[string]interface{}{ // ✅ REQUIRED — defines each argument
            "user_id":     map[string]string{"type": "string"},
            "destination": map[string]string{"type": "string"},
            "days":        map[string]string{"type": "number"},
            "budget":      map[string]string{"type": "number"},
        },
        "required": []string{"user_id", "destination"}, // ⚠️ OPTIONAL but strongly recommended
    },
},
```

---

### Why each key matters:

| Key              | Required?   | Purpose                                                                                  |
| ---------------- | ----------- | ---------------------------------------------------------------------------------------- |
| `name`           | ✅ Yes      | Unique ID for the tool. Used in `tools/call` → `params.name`                             |
| `description`    | ✅ Yes      | The **LLM reads this** to decide which tool to call. Bad description = wrong tool picked |
| `inputSchema`    | ✅ Yes      | JSON Schema telling the LLM **what arguments to generate**                               |
| `type: "object"` | ✅ Yes      | MCP always expects object-type schema                                                    |
| `properties`     | ✅ Yes      | Defines argument names and their types                                                   |
| `required`       | ⚠️ Optional | If missing, LLM may skip important args. Always include it for mandatory fields          |

---

**Key insight:** MCP is **schema-driven** — the LLM uses `description` + `inputSchema` to autonomously decide:

1. _Should I call this tool?_ (from `description`)
2. _What arguments do I pass?_ (from `inputSchema` + `properties`)

If any of these are wrong or missing, the AI will either call the wrong tool or generate bad arguments.
