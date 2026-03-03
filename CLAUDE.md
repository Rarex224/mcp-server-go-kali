# backless-mcp Implementation Details

## Overview
backless-mcp is a Go-based Model Context Protocol (MCP) server running in a multi-stage Docker environment (compiling via golang, running in kalilinux/kali-rolling).

## Guidelines for MCP Tool Development
- The application uses `github.com/mark3labs/mcp-go`.
- All inputs from the user are passed through `sanitizeInput` to strip shell metacharacters and restrict parameters to safe chars (`[a-zA-Z0-9.\-:/=?&]`). This avoids command injection via `os/exec`.
- Tools output directly to standard output/error, which the server captures and returns as strings `mcp.NewToolResultText()`.
- If no output is produced by a tool, it returns a success message rather than an empty string to keep the LLM informed.

## Running Tests / Validation
Since the application relies on external Kali tools, it is best tested by building and running the Docker image:
```bash
docker build -t backless-mcp-server .
docker run -i --rm backless-mcp-server
```
You can send JSON-RPC stdio commands like `{"jsonrpc":"2.0","method":"tools/list","id":1}` via stdin.

## Security
- Always ensure target URLs or IPs belong to an authorized environment.
- Tools are executed via Go's `os/exec` instead of `sh -c` to mitigate risk.
- Only safe string arguments are passed to the tools.
