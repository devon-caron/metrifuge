# Build stage
FROM golang:1.24-bullseye AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application with static linking
# CGO_ENABLED=0 ensures static binary
# -ldflags='-w -s' strips debug info to reduce size
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# Final stage - Debian slim with debugging tools
FROM debian:12-slim

# Install essential tools and utilities
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    wget \
    vim \
    htop \
    procps \
    net-tools \
    strace \
    tcpdump \
    dnsutils \
    iputils-ping \
    tree \
    less \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Create a non-root user
RUN groupadd -r -g 1001 appgroup && \
    useradd -r -u 1001 -g appgroup -m -s /bin/bash appuser

# Set the working directory
WORKDIR /app

# Template environment variables
ENV MF_LOG_LEVEL=debug
ENV MF_LOG_REPORTCALLER_STATUS=true
ENV MF_RUNNING_IN_K8S=true

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown appuser:appgroup /app/main

# Switch to non-root user
USER appuser

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]