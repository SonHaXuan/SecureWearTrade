# JEDI API - Containerized Encryption & Security Testing Platform

🚀 **JEDI (Java-informed Digital Encryption Infrastructure)** is a comprehensive API platform for encryption, decryption, and security testing, now fully containerized with Docker support.

## 📋 Table of Contents

- [Quick Start](#quick-start)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the Container](#running-the-container)
- [API Endpoints](#api-endpoints)
- [Usage Examples](#usage-examples)
- [Performance Metrics](#performance-metrics)
- [Troubleshooting](#troubleshooting)
- [Architecture](#architecture)

## 🚀 Quick Start

```bash
# Build the JEDI container
docker build -t jedi-encrypted .

# Run the container
docker run -d -p 8081:8080 --name jedi-api jedi-encrypted

# Test the API
curl http://localhost:8081/health
```

## 📋 Prerequisites

- **Docker** (version 20.10 or higher)
- **Docker Compose** (optional)
- **curl** or any HTTP client for API testing

## 🔧 Installation

### Option 1: Using Pre-built Container
```bash
# If you have the pre-built image
docker run -d -p 8081:8080 --name jedi-api jedi-encrypted
```

### Option 2: Build from Source
```bash
# Clone the repository (if applicable)
cd go-jedi

# Build the Docker image
docker build -t jedi-encrypted .

# Verify the image was built
docker images | grep jedi-encrypted
```

## 🐳 Running the Container

### Basic Run
```bash
docker run -d -p 8081:8080 --name jedi-api jedi-encrypted
```

### Advanced Run with Custom Settings
```bash
docker run -d \
  -p 8081:8080 \
  --name jedi-api \
  --restart unless-stopped \
  --memory=512m \
  --cpus=0.5 \
  jedi-encrypted
```

### Using Docker Compose
Create a `docker-compose.yml` file:

```yaml
version: '3.8'
services:
  jedi-api:
    build: .
    ports:
      - "8081:8080"
    container_name: jedi-api
    restart: unless-stopped
    environment:
      - GIN_MODE=release
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

Run with:
```bash
docker-compose up -d
```

## 📡 API Endpoints

### Base URL
```
http://localhost:8081
```

### Available Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|-----------------|
| GET | `/` | API Information | None |
| GET | `/health` | Health Check | None |
| POST | `/encrypt` | Encrypt Message | None |
| POST | `/decrypt` | Decrypt Message | None |

## 🔐 API Usage Examples

### 1. API Information
```bash
curl -X GET http://localhost:8081/ | jq .
```

**Response:**
```json
{
  "service": "JEDI Enhanced API",
  "status": "running",
  "timestamp": "2025-10-30T10:47:19.515438Z",
  "message": "JEDI container with encryption is working successfully!",
  "endpoints": ["/", "/health", "/encrypt", "/decrypt"]
}
```

### 2. Health Check
```bash
curl -X GET http://localhost:8081/health | jq .
```

**Response:**
```json
{
  "service": "JEDI Enhanced Health Check",
  "status": "healthy",
  "timestamp": "2025-10-30T10:47:19.515438Z",
  "message": "All systems operational including encryption!",
  "endpoints": []
}
```

### 3. Encrypt Message
```bash
curl -X POST http://localhost:8081/encrypt \
  -H "Content-Type: application/json" \
  -d '{
    "uri": "hospital/patient123/record",
    "message": "This is a confidential patient record with sensitive health data"
  }' | jq .
```

**Response:**
```json
{
  "success": true,
  "data": "Aoew3zDIIgHryh4weCNM0hqx0hxKS4VZfoOR9sX+/NJaHOYdQjLwRpNrwv+7aRBYe2LdZyJvc3B4qy7Fo69lzNZsuxU/3mcA+I51XFwFD30=",
  "executionTime": 372,
  "memoryUsage": 52428800,
  "cpuPercentage": 15.5,
  "powerUsage": 2.5,
  "energyConsumptionJoules": 0.00093
}
```

### 4. Decrypt Message
```bash
curl -X POST http://localhost:8081/decrypt \
  -H "Content-Type: application/json" \
  -d '{
    "uri": "hospital/patient123/record",
    "encryptedMessage": "Aoew3zDIIgHryh4weCNM0hqx0hxKS4VZfoOR9sX+/NJaHOYdQjLwRpNrwv+7aRBYe2LdZyJvc3B4qy7Fo69lzNZsuxU/3mcA+I51XFwFD30="
  }' | jq .
```

**Response:**
```json
{
  "success": true,
  "data": "This is a confidential patient record with sensitive health data",
  "executionTime": 219,
  "memoryUsage": 52428800,
  "cpuPercentage": 15.5,
  "powerUsage": 2.5,
  "energyConsumptionJoules": 0.00055
}
```

## 📊 Performance Metrics

The JEDI API provides detailed performance metrics for each encryption/decryption operation:

| Metric | Description | Example Value |
|--------|-------------|---------------|
| `executionTime` | Time taken in microseconds | 372 |
| `memoryUsage` | Memory usage in bytes | 52428800 |
| `cpuPercentage` | CPU utilization percentage | 15.5 |
| `powerUsage` | Estimated power consumption in watts | 2.5 |
| `energyConsumptionJoules` | Energy consumed in joules | 0.00093 |

### Performance Benchmarks

- **Encryption Speed**: 300-600 microseconds
- **Decryption Speed**: 200-300 microseconds
- **Memory Usage**: ~50MB constant
- **Power Efficiency**: ~0.001 Joules per operation

## 🔧 Configuration

### Environment Variables
```bash
# Set production mode
docker run -d -p 8081:8080 -e GIN_MODE=release --name jedi-api jedi-encrypted
```

### Container Resource Limits
```bash
docker run -d \
  -p 8081:8080 \
  --memory=512m \
  --cpus=0.5 \
  --name jedi-api \
  jedi-encrypted
```

## 🛠️ Management Commands

### Check Container Status
```bash
docker ps | grep jedi-api
```

### View Logs
```bash
docker logs jedi-api
```

### Follow Logs (Real-time)
```bash
docker logs -f jedi-api
```

### Stop Container
```bash
docker stop jedi-api
```

### Remove Container
```bash
docker rm jedi-api
```

### Restart Container
```bash
docker restart jedi-api
```

### View Container Statistics
```bash
docker stats jedi-api
```

## 🔍 Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Error: "port is already allocated"
# Solution: Use a different port or stop the conflicting container
docker run -d -p 8082:8080 --name jedi-api jedi-encrypted
```

#### 2. Container Not Starting
```bash
# Check logs for errors
docker logs jedi-api

# Common solutions:
# - Ensure port 8081 is available
# - Check Docker daemon is running
# - Verify image was built successfully
```

#### 3. API Not Responding
```bash
# Check if container is running
docker ps | grep jedi-api

# Check container health
curl http://localhost:8081/health

# Restart if needed
docker restart jedi-api
```

#### 4. Memory Issues
```bash
# Increase memory allocation
docker run -d -p 8081:8080 --memory=1g --name jedi-api jedi-encrypted
```

### Debug Mode
```bash
# Run container in debug mode
docker run -it --rm -p 8081:8080 jedi-encrypted

# Or run directly without Docker
go run enhanced_main.go
```

## 📈 Monitoring

### Health Monitoring
```bash
# Continuous health check
watch -n 5 'curl -s http://localhost:8081/health | jq .'
```

### Performance Monitoring
```bash
# Monitor container resources
docker stats --no-stream jedi-api

# Monitor API response time
curl -w "@curl-format.txt" -s -o /dev/null http://localhost:8081/health
```

### Log Analysis
```bash
# Count API requests
docker logs jedi-api | grep "POST" | wc -l

# Monitor errors
docker logs jedi-api | grep -i error
```

## 🏗️ Architecture

### Container Architecture
```
┌─────────────────────────────────────┐
│           Docker Container           │
│  ┌─────────────────────────────────┐ │
│  │        Go Application           │ │
│  │  ┌─────────────┬─────────────┐  │ │
│  │  │ HTTP Server │  Encryption  │  │ │
│  │  │   (Port     │   (AES-256   │  │ │
│  │  │   8080)     │   CFB Mode)  │  │ │
│  │  └─────────────┴─────────────┘  │ │
│  └─────────────────────────────────┘ │
└─────────────────────────────────────┘
```

### Technology Stack
- **Language**: Go 1.20.5
- **Web Framework**: Standard library `net/http`
- **Encryption**: AES-256-CFB
- **Container**: Docker with Ubuntu 24.04
- **Architecture**: ARM64/AMD64 compatible

### Security Features
- **AES-256 Encryption**: Military-grade encryption
- **Context-Based**: URI-specific encryption contexts
- **Base64 Encoding**: Safe data transmission
- **Performance Monitoring**: Resource usage tracking
- **Energy Analysis**: Power consumption metrics

## 📝 Development

### Local Development
```bash
# Install Go 1.20.5+
# Clone the repository
cd go-jedi

# Run locally
go run enhanced_main.go

# Build binary
go build -o jedi-api enhanced_main.go

# Run binary
./jedi-api
```

### Testing
```bash
# Run health check
curl http://localhost:8080/health

# Test encryption
curl -X POST http://localhost:8080/encrypt \
  -H "Content-Type: application/json" \
  -d '{"uri":"test","message":"hello world"}'

# Test decryption
curl -X POST http://localhost:8080/decrypt \
  -H "Content-Type: application/json" \
  -d '{"uri":"test","encryptedMessage":"<encrypted_data>"}'
```

### Building Different Versions
```bash
# Minimal version (no encryption)
docker build -t jedi-minimal -f Dockerfile.minimal .

# Full version with all features
docker build -t jedi-full .

# Custom version
docker build -t jedi-custom -f Dockerfile.custom .
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE.txt file for details.

## 📞 Support

For support and questions:
- Check the troubleshooting section
- Review the API documentation
- Create an issue in the repository
- Check container logs for error messages

---

## 🎯 Quick Reference Commands

```bash
# Build and run
docker build -t jedi-encrypted . && docker run -d -p 8081:8080 --name jedi-api jedi-encrypted

# Test health
curl http://localhost:8081/health

# Test encryption
curl -X POST http://localhost:8081/encrypt -H "Content-Type: application/json" -d '{"uri":"test","message":"hello"}'

# View logs
docker logs jedi-api

# Stop and clean
docker stop jedi-api && docker rm jedi-api
```

**🚀 JEDI API is ready for production use!**