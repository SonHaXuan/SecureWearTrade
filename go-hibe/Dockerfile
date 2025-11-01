# Use the official Ubuntu image from the Docker Hub
FROM ubuntu:24.04

# Set the Go version as an environment variable
ENV GO_VERSION=1.20.5

# Install necessary packages
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    git \
    build-essential \
    clang \
    lsb-release && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Download and install Go (supporting multiple architectures)
RUN ARCH=$(dpkg --print-architecture) && \
    if [ "$ARCH" = "arm64" ]; then \
        wget https://golang.org/dl/go$GO_VERSION.linux-arm64.tar.gz && \
        tar -C /usr/local -xzf go$GO_VERSION.linux-arm64.tar.gz && \
        rm go$GO_VERSION.linux-arm64.tar.gz; \
    else \
        wget https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz && \
        tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz && \
        rm go$GO_VERSION.linux-amd64.tar.gz; \
    fi

# Set Go paths
ENV PATH=/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV GOPATH=/go
ENV GO111MODULE=on
ENV GOCACHE=/root/.cache/go-build
ENV PATH=$GOPATH/bin:$PATH

# Verify Go installation
RUN go version

# Default command to run when container starts
# CMD ["bash"]

WORKDIR /app

# Copy entire source code first (needed for local package replacements)
COPY . .

# No external dependencies needed for minimal version
# Skip dependency resolution

# Build and run the application
CMD ["go", "run", "enhanced_main.go"]
