# Testing backless-mcp

This project includes automated testing workflows to verify the functionality of the MCP server and its underlying tools.

## Prerequisites
To run these tests, you must have installed:
1. **Docker Desktop** (Required for the Kali Linux environment)
2. **Go (1.22+)** (Optional, only needed if you want to run unit tests natively outside Docker)

## 1. Automated Integration Tests (Requires Docker)
We have provided an automated bash script that builds the Docker image and simulates the JSON-RPC communication that an MCP client (like Claude Desktop) would perform.

Run the test script from your terminal:
```bash
cd mcp/backless-mcp
./run_tests.sh
```

**What it does:**
- Verifies Docker is installed.
- Builds the `backless-mcp-server` image.
- Sends an `initialize` request to the MCP server.
- Sends a `tools/list` request to verify all Kali tools are exposed.
- Sends a `tools/call` request to execute a safe command (`whois_lookup`) and verifies the output.

## 2. Unit Tests (Requires Go)
We have included a `main_test.go` file which tests the core Go logic, specifically ensuring that the `sanitizeInput` function properly strips out malicious characters to prevent command injection vulnerabilities.

Run the unit tests:
```bash
cd mcp/backless-mcp
go test -v ./...
```
