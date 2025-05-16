package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/semgrep-cli-mcp/internal/mcpx"
	"github.com/semgrep-cli-mcp/internal/semgrep"
)

type ListConfigsTool struct {
	ConfigDir string
}

func (t *ListConfigsTool) ServerTool() mcpserver.ServerTool {
	return mcpserver.ServerTool{
		Tool: mcp.NewTool("list_configs",
			mcp.WithDescription("List available Semgrep configurations"),
		),
		Handler: t.Handle,
	}
}

func (t *ListConfigsTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	configs, err := semgrep.ReadConfigs(t.ConfigDir)
	if err != nil {
		return mcpx.NewToolResultErrorf("%v", err), nil
	}
	return mcpx.NewToolResultJSON(configs), nil
}
