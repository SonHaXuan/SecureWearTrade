# JEDI API Detailed Documentation

## Overview

The JEDI API provides RESTful endpoints for encryption, decryption, and system monitoring. This document provides detailed technical specifications for each endpoint.

## Base URL
```
http://localhost:8081
```

## Content-Type
All requests must include:
```
Content-Type: application/json
```

## Response Format
All responses are in JSON format with the following structure:
```json
{
  "success": boolean,
  "data": any,
  "error": string,
  "timestamp": string,
  "executionTime": number,
  "memoryUsage": number,
  "cpuPercentage": number,
  "powerUsage": number,
  "energyConsumptionJoules": number
}
```

## Endpoints

### 1. Get API Information
Retrieve basic information about the JEDI API service.

**Endpoint**: `GET /`

**Description**: Returns service status, available endpoints, and general information.

**Request**: No parameters required

**Example Request**:
```bash
curl -X GET http://localhost:8081/
```

**Response Schema**:
```json
{
  "service": "string",
  "status": "running|stopped|error",
  "timestamp": "ISO 8601 datetime",
  "message": "string",
  "endpoints": ["string"]
}
```

**Example Response**:
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
Check the health status of the JEDI API service.

**Endpoint**: `GET /health`

**Description**: Returns detailed health information about the service.

**Request**: No parameters required

**Example Request**:
```bash
curl -X GET http://localhost:8081/health
```

**Response Schema**:
```json
{
  "service": "string",
  "status": "healthy|unhealthy",
  "timestamp": "ISO 8601 datetime",
  "message": "string",
  "endpoints": []
}
```

**Example Response**:
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
Encrypt a message using AES-256-CFB encryption.

**Endpoint**: `POST /encrypt`

**Description**: Encrypts a provided message and returns the encrypted data with performance metrics.

**Request Schema**:
```json
{
  "uri": "string",
  "message": "string"
}
```

**Request Parameters**:
- `uri` (required, string): Context URI for the encryption (e.g., "hospital/patient123/record")
- `message` (required, string): Message to be encrypted

**Example Request**:
```bash
curl -X POST http://localhost:8081/encrypt \
  -H "Content-Type: application/json" \
  -d '{
    "uri": "hospital/patient123/record",
    "message": "This is a confidential patient record with sensitive health data"
  }'
```

**Response Schema**:
```json
{
  "success": boolean,
  "data": "string",
  "error": "string",
  "executionTime": number,
  "memoryUsage": number,
  "cpuPercentage": number,
  "powerUsage": number,
  "energyConsumptionJoules": number
}
```

**Response Parameters**:
- `success`: Indicates if encryption was successful
- `data`: Base64-encoded encrypted data
- `error`: Error message if encryption failed
- `executionTime`: Time taken in microseconds
- `memoryUsage`: Memory usage in bytes
- `cpuPercentage`: CPU utilization percentage
- `powerUsage`: Estimated power consumption in watts
- `energyConsumptionJoules`: Energy consumed in joules

**Example Response**:
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
Decrypt a previously encrypted message.

**Endpoint**: `POST /decrypt`

**Description**: Decrypts a Base64-encoded encrypted message and returns the original data with performance metrics.

**Request Schema**:
```json
{
  "uri": "string",
  "encryptedMessage": "string"
}
```

**Request Parameters**:
- `uri` (required, string): Context URI that was used for encryption
- `encryptedMessage` (required, string): Base64-encoded encrypted data

**Example Request**:
```bash
curl -X POST http://localhost:8081/decrypt \
  -H "Content-Type: application/json" \
  -d '{
    "uri": "hospital/patient123/record",
    "encryptedMessage": "Aoew3zDIIgHryh4weCNM0hqx0hxKS4VZfoOR9sX+/NJaHOYdQjLwRpNrwv+7aRBYe2LdZyJvc3B4qy7Fo69lzNZsuxU/3mcA+I51XFwFD30="
  }'
```

**Response Schema**:
```json
{
  "success": boolean,
  "data": "string",
  "error": "string",
  "executionTime": number,
  "memoryUsage": number,
  "cpuPercentage": number,
  "powerUsage": number,
  "energyConsumptionJoules": number
}
```

**Response Parameters**:
- `success`: Indicates if decryption was successful
- `data`: Decrypted original message
- `error`: Error message if decryption failed
- `executionTime`: Time taken in microseconds
- `memoryUsage`: Memory usage in bytes
- `cpuPercentage`: CPU utilization percentage
- `powerUsage`: Estimated power consumption in watts
- `energyConsumptionJoules`: Energy consumed in joules

**Example Response**:
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

## Error Handling

### HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request |
| 405 | Method Not Allowed |
| 500 | Internal Server Error |

### Error Response Format

```json
{
  "success": false,
  "error": "Error description",
  "timestamp": "ISO 8601 datetime"
}
```

### Common Errors

#### Invalid JSON
```json
{
  "success": false,
  "error": "Invalid JSON: invalid character '}' looking for beginning of object key string"
}
```

#### Missing Parameters
```json
{
  "success": false,
  "error": "Invalid JSON: json: unknown field \"invalid_field\""
}
```

#### Method Not Allowed
```json
{
  "success": false,
  "error": "Method not allowed. Use POST."
}
```

#### Encryption/Decryption Errors
```json
{
  "success": false,
  "error": "Encryption failed: invalid message length"
}
```

## Performance Metrics

### Metrics Description

| Metric | Unit | Description |
|--------|------|-------------|
| executionTime | microseconds | Time taken to process the request |
| memoryUsage | bytes | Memory allocated during processing |
| cpuPercentage | percent | CPU utilization during processing |
| powerUsage | watts | Estimated power consumption |
| energyConsumptionJoules | joules | Total energy consumed |

### Typical Performance Ranges

| Operation | Min Time | Max Time | Average Time |
|-----------|----------|----------|--------------|
| Encrypt | 300μs | 600μs | 400μs |
| Decrypt | 200μs | 300μs | 250μs |
| Health Check | 10μs | 50μs | 20μs |

## Security Considerations

### Encryption Algorithm
- **Algorithm**: AES-256-CFB
- **Key Derivation**: SHA-256 hash of master key
- **IV Generation**: Cryptographically secure random IV per encryption
- **Encoding**: Base64 for safe transmission

### Data Protection
- Messages are encrypted in transit
- Each URI context provides different encryption parameters
- No plaintext data is stored in logs
- Memory is cleared after processing

### Best Practices
1. Use HTTPS in production environments
2. Implement proper API authentication
3. Validate input data before encryption
4. Monitor for unusual activity patterns
5. Regularly rotate encryption keys

## Rate Limiting

Currently, the API does not implement rate limiting. In production, consider implementing:
- Request rate limiting per IP
- Maximum message size limits
- Concurrent request limits

## Versioning

API versioning follows semantic versioning:
- Current version: v1.0.0
- Breaking changes will increment major version
- New features will increment minor version
- Bug fixes will increment patch version

## Testing

### Unit Tests
```bash
go test ./...
```

### Integration Tests
```bash
# Test all endpoints
curl -X GET http://localhost:8081/health

# Test encryption cycle
MESSAGE="test message"
URI="test/context"

# Encrypt
ENCRYPTED=$(curl -s -X POST http://localhost:8081/encrypt \
  -H "Content-Type: application/json" \
  -d "{\"uri\":\"$URI\",\"message\":\"$MESSAGE\"}" | jq -r '.data')

# Decrypt
curl -s -X POST http://localhost:8081/decrypt \
  -H "Content-Type: application/json" \
  -d "{\"uri\":\"$URI\",\"encryptedMessage\":\"$ENCRYPTED\"}"
```

### Load Testing
```bash
# Install Apache Bench
ab -n 1000 -c 10 http://localhost:8081/health

# Test encryption endpoint under load
ab -n 100 -c 5 -p encrypt_request.json -T application/json http://localhost:8081/encrypt
```

## Monitoring and Logging

### Log Format
```
🚀 JEDI Enhanced API Server Starting...
✅ JEDI Enhanced API Server started successfully!
🌐 Server running on http://localhost:8080
📊 Health check: http://localhost:8080/health
🔐 Encrypt API: http://localhost:8080/encrypt
🔓 Decrypt API: http://localhost:8080/decrypt
🐳 Docker container with encryption is running!
```

### Health Monitoring
Regular health checks should be performed:
```bash
# Check every 30 seconds
watch -n 30 'curl -s http://localhost:8081/health | jq .status'
```

### Performance Monitoring
Monitor key metrics:
- Response times
- Error rates
- Memory usage
- CPU utilization
- Power consumption

## SDK Examples

### Python Example
```python
import requests
import json

BASE_URL = "http://localhost:8081"

def encrypt_message(uri, message):
    response = requests.post(
        f"{BASE_URL}/encrypt",
        json={"uri": uri, "message": message}
    )
    return response.json()

def decrypt_message(uri, encrypted_message):
    response = requests.post(
        f"{BASE_URL}/decrypt",
        json={"uri": uri, "encryptedMessage": encrypted_message}
    )
    return response.json()

# Usage
result = encrypt_message("test/uri", "Hello, World!")
print(f"Encrypted: {result['data']}")

decrypted = decrypt_message("test/uri", result['data'])
print(f"Decrypted: {decrypted['data']}")
```

### JavaScript Example
```javascript
const BASE_URL = 'http://localhost:8081';

async function encryptMessage(uri, message) {
    const response = await fetch(`${BASE_URL}/encrypt`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ uri, message }),
    });
    return await response.json();
}

async function decryptMessage(uri, encryptedMessage) {
    const response = await fetch(`${BASE_URL}/decrypt`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ uri, encryptedMessage }),
    });
    return await response.json();
}

// Usage
encryptMessage('test/uri', 'Hello, World!')
    .then(result => {
        console.log('Encrypted:', result.data);
        return decryptMessage('test/uri', result.data);
    })
    .then(decrypted => {
        console.log('Decrypted:', decrypted.data);
    });
```

---

This API documentation provides comprehensive information for developers to integrate with the JEDI encryption service effectively.