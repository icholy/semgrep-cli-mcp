package main

import (
	"flag"
	"fmt"
	"os"

	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/semgrep-cli-mcp/tools"
)

func main() {
	var configDir string
	flag.StringVar(&configDir, "configs", "./semgrep", "Directory containing Semgrep config files")
	flag.Parse()

	server := mcpserver.NewMCPServer(
		"Semgrep CLI MCP",
		"1.0.0",
		mcpserver.WithToolCapabilities(false),
		mcpserver.WithRecovery(),
	)

	listConfigsTool := &tools.ListConfigsTool{ConfigDir: configDir}
	scanTool := &tools.ScanTool{ConfigDir: configDir}

	server.AddTools(
		listConfigsTool.ServerTool(),
		scanTool.ServerTool(),
	)

	if err := mcpserver.ServeStdio(server); err != nil {
		fmt.Printf("mcpserver error: %v\n", err)
		os.Exit(1)
	}
}
