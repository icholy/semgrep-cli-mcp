package tools

import (
	"context"
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/semgrep-cli-mcp/internal/mcpx"
	"github.com/semgrep-cli-mcp/internal/semgrep"
)

type ScanTool struct {
	ConfigDir string
}

func (t *ScanTool) ServerTool() mcpserver.ServerTool {
	return mcpserver.ServerTool{
		Tool: mcp.NewTool("scan",
			mcp.WithDescription("Run Semgrep scan with specified configuration. The results contain the exact file path and line numbers."),
			mcp.WithString("directory",
				mcp.DefaultString("."),
				mcp.Description("Directory to scan"),
			),
			mcp.WithString("config",
				mcp.Required(),
				mcp.Description("Name of the configuration to use. Note use the list_configs tool to enumerate them"),
			),
		),
		Handler: t.Handle,
	}
}

type SemgrepResult struct {
	File  string `json:"file"`
	Line  int    `json:"line"`
	Lines string `json:"lines"`
}

func (t *ScanTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	directory, ok := mcpx.GetParamArgument[string](request, "directory")
	if !ok {
		directory = "."
	}
	config, ok := mcpx.GetParamArgument[string](request, "config")
	if !ok {
		return mcpx.NewToolResultErrorf("Missing config argument"), nil
	}
	output, err := semgrep.Scan(semgrep.ScanOptions{
		Dir:        directory,
		ConfigPath: filepath.Join(t.ConfigDir, config),
	})
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && len(exitErr.Stderr) > 0 {
			return mcpx.NewToolResultErrorf("Scan failed: %v\nStderr: %s", err, string(exitErr.Stderr)), nil
		}
		return mcpx.NewToolResultErrorf("Scan failed: %v", err), nil
	}
	var results []SemgrepResult
	for _, r := range output.Results {
		lines, err := semgrep.ReadLines(r, semgrep.ReadLinesOptions{
			Dir:    directory,
			Extend: true,
		})
		if err != nil {
			return mcpx.NewToolResultErrorf("Failed to read lines: %v", err), nil
		}
		results = append(results, SemgrepResult{
			File:  r.Path,
			Line:  r.Start.Line,
			Lines: lines,
		})
	}
	return mcpx.NewToolResultJSON(results), nil
}
