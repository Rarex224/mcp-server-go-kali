package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// sanitizeInput performs basic sanitization to prevent command injection
func sanitizeInput(input string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9.\-:/=?&_,@%]`)
	return reg.ReplaceAllString(input, "")
}

func runToolCommand(name string, args ...string) (*mcp.CallToolResult, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	
	result := string(out)
	if err != nil {
		result = fmt.Sprintf("Error executing %s: %v\nOutput:\n%s", name, err, result)
		return mcp.NewToolResultText(result), nil
	}
	
	if result == "" {
		result = "✅ Command executed successfully with no output."
	}
	
	return mcp.NewToolResultText(result), nil
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer("backless-mcp", "1.1.0")

	// 1. Nmap Tool
	nmapTool := mcp.NewTool("nmap_scan",
		mcp.WithDescription("Run Nmap scan against a target IP or domain"),
		mcp.WithString("target", 
			mcp.Required(), 
			mcp.Description("Target IP or domain to scan"),
		),
		mcp.WithString("flags", 
			mcp.Description("Optional Nmap flags (e.g., -p 80,443 -sV). Restricted characters removed."),
		),
	)
	s.AddTool(nmapTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		target, _ := request.RequireString("target")
		target = sanitizeInput(target)
		flagsStr := request.GetString("flags", "")
		
		args := []string{}
		if flagsStr != "" {
			flags := strings.Fields(flagsStr)
			for _, f := range flags {
				args = append(args, sanitizeInput(f))
			}
		}
		args = append(args, target)
		return runToolCommand("nmap", args...)
	})

	// 2. Nikto Tool
	niktoTool := mcp.NewTool("nikto_scan",
		mcp.WithDescription("Run Nikto web scanner against a target URL"),
		mcp.WithString("url", 
			mcp.Required(), 
			mcp.Description("Target URL (e.g., http://example.com)"),
		),
	)
	s.AddTool(niktoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, _ := request.RequireString("url")
		target := sanitizeInput(url)
		args := []string{"-h", target}
		return runToolCommand("nikto", args...)
	})

	// 3. SQLMap Tool
	sqlmapTool := mcp.NewTool("sqlmap_scan",
		mcp.WithDescription("Run SQLMap against a target URL for basic injection testing"),
		mcp.WithString("url", 
			mcp.Required(), 
			mcp.Description("Target URL with parameters (e.g., http://example.com/page?id=1)"),
		),
	)
	s.AddTool(sqlmapTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, _ := request.RequireString("url")
		target := sanitizeInput(url)
		args := []string{"-u", target, "--batch", "--random-agent"}
		return runToolCommand("sqlmap", args...)
	})

	// 4. WPScan Tool
	wpscanTool := mcp.NewTool("wpscan",
		mcp.WithDescription("Run WPScan against a target WordPress site"),
		mcp.WithString("url", 
			mcp.Required(), 
			mcp.Description("Target URL"),
		),
	)
	s.AddTool(wpscanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, _ := request.RequireString("url")
		target := sanitizeInput(url)
		args := []string{"--url", target, "--no-banner", "--random-user-agent"}
		return runToolCommand("wpscan", args...)
	})

	// 5. Dirb Tool
	dirbTool := mcp.NewTool("dirb_scan",
		mcp.WithDescription("Run Dirb directory brute-forcer against a target URL"),
		mcp.WithString("url", 
			mcp.Required(), 
			mcp.Description("Target URL"),
		),
	)
	s.AddTool(dirbTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, _ := request.RequireString("url")
		target := sanitizeInput(url)
		args := []string{target, "-r"} // non-recursive by default
		return runToolCommand("dirb", args...)
	})

	// 6. Searchsploit Tool
	searchsploitTool := mcp.NewTool("searchsploit",
		mcp.WithDescription("Search for exploits using exploitdb (searchsploit)"),
		mcp.WithString("query", 
			mcp.Required(), 
			mcp.Description("Search query (e.g., apache 2.4)"),
		),
	)
	s.AddTool(searchsploitTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := request.RequireString("query")
		q := sanitizeInput(query)
		args := []string{q}
		return runToolCommand("searchsploit", args...)
	})

	// 7. Gobuster Tool
	gobusterTool := mcp.NewTool("gobuster_dir",
		mcp.WithDescription("Run Gobuster directory brute-forcer against a target URL"),
		mcp.WithString("url", 
			mcp.Required(), 
			mcp.Description("Target URL"),
		),
	)
	s.AddTool(gobusterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, _ := request.RequireString("url")
		target := sanitizeInput(url)
		args := []string{"dir", "-u", target, "-w", "/usr/share/wordlists/dirb/common.txt", "-q"}
		return runToolCommand("gobuster", args...)
	})

	// 8. Whois Tool
	whoisTool := mcp.NewTool("whois_lookup",
		mcp.WithDescription("Perform a WHOIS lookup for a domain"),
		mcp.WithString("domain", 
			mcp.Required(), 
			mcp.Description("Domain name (e.g., example.com)"),
		),
	)
	s.AddTool(whoisTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		domain, _ := request.RequireString("domain")
		d := sanitizeInput(domain)
		args := []string{d}
		return runToolCommand("whois", args...)
	})

	// 9. Dnsenum Tool
	dnsenumTool := mcp.NewTool("dnsenum",
		mcp.WithDescription("Perform DNS enumeration for a domain using dnsenum"),
		mcp.WithString("domain", 
			mcp.Required(), 
			mcp.Description("Domain name (e.g., example.com)"),
		),
	)
	s.AddTool(dnsenumTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		domain, _ := request.RequireString("domain")
		d := sanitizeInput(domain)
		args := []string{"--noreverse", d}
		return runToolCommand("dnsenum", args...)
	})

	// 10. SMBMap Tool
	smbmapTool := mcp.NewTool("smbmap",
		mcp.WithDescription("Enumerate SMB shares on a target IP"),
		mcp.WithString("host", 
			mcp.Required(), 
			mcp.Description("Target IP address"),
		),
	)
	s.AddTool(smbmapTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host, _ := request.RequireString("host")
		h := sanitizeInput(host)
		args := []string{"-H", h}
		return runToolCommand("smbmap", args...)
	})

	// 11. Enum4linux Tool
	enum4linuxTool := mcp.NewTool("enum4linux",
		mcp.WithDescription("Enumerate information from Windows and Samba systems"),
		mcp.WithString("host", 
			mcp.Required(), 
			mcp.Description("Target IP address"),
		),
	)
	s.AddTool(enum4linuxTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		host, _ := request.RequireString("host")
		h := sanitizeInput(host)
		args := []string{"-a", h}
		return runToolCommand("enum4linux", args...)
	})

	// 12. WhatWeb Tool
	whatwebTool := mcp.NewTool("whatweb",
		mcp.WithDescription("Identify website technologies using WhatWeb"),
		mcp.WithString("url", 
			mcp.Required(), 
			mcp.Description("Target URL"),
		),
	)
	s.AddTool(whatwebTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, _ := request.RequireString("url")
		u := sanitizeInput(url)
		args := []string{u}
		return runToolCommand("whatweb", args...)
	})

	// 13. TheHarvester Tool
	theharvesterTool := mcp.NewTool("theharvester",
		mcp.WithDescription("Gather OSINT emails, subdomains, and IPs for a domain"),
		mcp.WithString("domain", 
			mcp.Required(), 
			mcp.Description("Target domain"),
		),
	)
	s.AddTool(theharvesterTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		domain, _ := request.RequireString("domain")
		d := sanitizeInput(domain)
		// Use a couple of basic sources that don't require API keys
		args := []string{"-d", d, "-b", "duckduckgo,crtsh", "-l", "100"}
		return runToolCommand("theHarvester", args...)
	})

	// Determine the port for HTTP/SSE transport, e.g., for Railway deployments
	port := os.Getenv("PORT")
	if port != "" {
		log.Printf("Starting backless-mcp server on SSE (HTTP) port %s", port)
		// Railway usually expects the server to listen on 0.0.0.0:PORT
		// We use NewSSEServer and Start it on the specified port.
		sseServer := server.NewSSEServer(s, server.WithBaseURL("http://0.0.0.0:"+port))
		if err := sseServer.Start(":" + port); err != nil {
			log.Fatalf("SSE Server error: %v", err)
		}
	} else {
		log.Println("Starting backless-mcp server on stdio (Local Claude Desktop)")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Stdio Server error: %v", err)
		}
	}
}
