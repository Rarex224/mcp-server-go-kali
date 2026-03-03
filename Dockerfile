# Stage 1: Build the Go binary
FROM golang:1.24-bookworm AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download || true

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o backless-mcp .

# Stage 2: Create the Kali Linux runtime image
FROM kalilinux/kali-rolling

ENV DEBIAN_FRONTEND=noninteractive

# Install required Kali tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    nmap nikto sqlmap wpscan dirb exploitdb \
    ca-certificates libcap2-bin \
    gobuster whois dnsenum smbmap enum4linux theharvester tcpdump netcat-traditional whatweb hydra john \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN useradd -m -u 1000 mcpuser
RUN setcap cap_net_raw,cap_net_admin,cap_net_bind_service+eip /usr/bin/nmap || true
RUN setcap cap_net_raw,cap_net_admin+eip /usr/sbin/tcpdump || true

WORKDIR /app
COPY --from=builder /app/backless-mcp .
RUN chown -R mcpuser:mcpuser /app

ENV PYTHONUNBUFFERED=1
USER mcpuser

CMD ["./backless-mcp"]
