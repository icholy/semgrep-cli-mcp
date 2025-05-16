# Local Semgrep CLI MCP

> An MCP server for using the local Semgrep CLI.

## Overview

This project provides an MCP server that exposes your local [Semgrep](https://semgrep.dev/) CLI as a set of tools, using a directory of Semgrep YAML configuration files.

## Tools

- **list_configs**
  Lists all available Semgrep configuration files in the configured directory.

- **scan**
  Runs a Semgrep scan using a specified configuration file from the config directory.

## Installation

Install the MCP server using Go:

```sh
go install github.com/icholy/semgrep-cli-mcp@latest
```

## Configuration

The server expects a directory containing Semgrep YAML config files (rulesets). Each config file should be a valid Semgrep rules YAML.

Example directory structure:

```
semgrep/
  python-security.yml
  js-best-practices.yml
  ...
```

## Cursor Config:

```json
{
  "mcpServers": {
    "semgrep-cli-mcp": {
      "command": "semgrep-cli-mcp",
      "args": [
        "--configs",
        "/path/to/your/semgrep/rules"
      ]
    }
  }
}
```