// Package mcpx provides utilities for MCP extensions.
package mcpx

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// NewToolResultErrorf returns a tool result containing a formatted error message.
func NewToolResultErrorf(format string, args ...interface{}) *mcp.CallToolResult {
	return mcp.NewToolResultText("Error: " + fmt.Sprintf(format, args...))
}

// NewToolResultJSON returns a tool result containing the marshaled JSON of v.
func NewToolResultJSON(v interface{}) *mcp.CallToolResult {
	data, err := json.Marshal(v)
	if err != nil {
		return NewToolResultErrorf("failed to marshal JSON: %v", err)
	}
	return mcp.NewToolResultText(string(data))
}

// GetParamArgument retrieves a parameter argument from the request and attempts to cast it to type T.
func GetParamArgument[T any](request mcp.CallToolRequest, key string) (T, bool) {
	if value, ok := request.Params.Arguments[key].(T); ok {
		return value, true
	}
	var zero T
	return zero, false
}
