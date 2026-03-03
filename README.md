# backless-mcp Server

A Model Context Protocol (MCP) server that provides access to various Kali Linux security tools via a Go-based MCP interface. It allows AI assistants to perform network scanning and vulnerability discovery using industry-standard CLI tools.

## Purpose

This MCP server provides a secure, containerized interface for AI assistants to run basic security scans using the top Kali Linux CLI tools.

## Features

### Current Implementation
- **`nmap_scan`** - Run Nmap scan against a target IP or domain
- **`nikto_scan`** - Run Nikto web scanner against a target URL
- **`sqlmap_scan`** - Run SQLMap against a target URL for basic injection testing
- **`wpscan`** - Run WPScan against a target WordPress site
- **`dirb_scan`** - Run Dirb directory brute-forcer against a target URL
- **`searchsploit`** - Search for exploits using exploitdb (searchsploit)
- **`gobuster_dir`** - Run Gobuster directory brute-forcer
- **`whois_lookup`** - Perform a WHOIS lookup for a domain
- **`dnsenum`** - Perform DNS enumeration for a domain
- **`smbmap`** - Enumerate SMB shares on a target IP
- **`enum4linux`** - Enumerate information from Windows and Samba systems
- **`whatweb`** - Identify website technologies
- **`theharvester`** - Gather OSINT emails, subdomains, and IPs for a domain

## Prerequisites

- Docker Desktop
- MCP Toolkit or an MCP client (e.g. Claude Desktop)

## Installation

### Step 1: Build Docker Image

```bash
cd mcp/backless-mcp
docker build -t backless-mcp-server .
```

### Step 2: Configure Claude Desktop

Find your Claude Desktop config file:
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

Add the server entry:
```json
{
  "mcpServers": {
    "backless-mcp": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "backless-mcp-server:latest"
      ]
    }
  }
}
```

### Step 3: Restart Claude Desktop
Restart your application. The tools should now be available.

## Usage Examples

In Claude Desktop, you can ask:
- "Run a basic nmap scan on example.com"
- "Use nikto to scan http://example.com"
- "Do a whois lookup on example.com"
- "Gather OSINT for example.com using theharvester"
- "Use gobuster to find hidden directories on http://example.com"

## Architecture

Claude Desktop -> Docker (stdio) -> backless-mcp (Go) -> Kali Linux CLI Tools

## Security Considerations

- All inputs are heavily sanitized to prevent command injection.
- Running as a non-root user (`mcpuser`) with specific Linux capabilities (`cap_net_raw`, etc.) enabled for network tooling.
- Be careful about scanning targets. **Always ensure you have authorization to scan the target environment.**

## License

MIT License
