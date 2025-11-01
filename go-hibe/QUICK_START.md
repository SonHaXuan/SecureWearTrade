# HIBE API - Quick Start Guide

🚀 **Get HIBE API running in under 2 minutes!**

## Prerequisites
- Docker installed on your system
- Terminal/command line access

## 🏃‍♂️ Quick Start (3 Commands)

```bash
# 1. Build the container
docker build -t hibe-encrypted .

# 2. Run the container
docker run -d -p 8081:8080 --name hibe-api hibe-encrypted

# 3. Test it's working
curl http://localhost:8081/health
```

Expected response:
```json
{
  "service": "HIBE Enhanced Health Check",
  "status": "healthy",
  "timestamp": "2025-10-30T10:47:19.515438Z",
  "message": "All systems operational including encryption!",
  "endpoints": []
}
```

## 🔐 Try Encryption (Copy & Paste)

```bash
# Encrypt a message
curl -X POST http://localhost:8081/encrypt \
  -H "Content-Type: application/json" \
  -d '{"uri":"test","message":"Hello HIBE!"}' | jq .
```

You'll get an encrypted response like:
```json
{
  "success": true,
  "data": "xP8kL2mN9qR7sT4vW1yZ3aB5cD8eF0gH=",
  "executionTime": 372,
  "memoryUsage": 52428800,
  "cpuPercentage": 15.5,
  "powerUsage": 2.5,
  "energyConsumptionJoules": 0.00093
}
```

## 🔓 Try Decryption

```bash
# Decrypt the message (use the encrypted data from above)
curl -X POST http://localhost:8081/decrypt \
  -H "Content-Type: application/json" \
  -d '{"uri":"test","encryptedMessage":"xP8kL2mN9qR7sT4vW1yZ3aB5cD8eF0gH="}' | jq .
```

You'll get back your original message:
```json
{
  "success": true,
  "data": "Hello HIBE!",
  "executionTime": 219,
  "memoryUsage": 52428800,
  "cpuPercentage": 15.5,
  "powerUsage": 2.5,
  "energyConsumptionJoules": 0.00055
}
```

## 🛠️ Common Commands

### Check if running
```bash
docker ps | grep hibe-api
```

### View logs
```bash
docker logs hibe-api
```

### Stop container
```bash
docker stop hibe-api
```

### Restart container
```bash
docker restart hibe-api
```

### Remove container
```bash
docker rm hibe-api
```

## 📡 All Available Endpoints

| Method | URL | Description |
|--------|-----|-------------|
| GET | `http://localhost:8081/` | API Info |
| GET | `http://localhost:8081/health` | Health Check |
| POST | `http://localhost:8081/encrypt` | Encrypt Message |
| POST | `http://localhost:8081/decrypt` | Decrypt Message |

## 🚨 Troubleshooting

### Port already in use?
```bash
# Try a different port
docker run -d -p 8082:8080 --name hibe-api hibe-encrypted
```

### Container not starting?
```bash
# Check for errors
docker logs hibe-api

# Try rebuilding
docker build -t hibe-encrypted . --no-cache
```

### API not responding?
```bash
# Check if container is running
docker ps | grep hibe-api

# Test basic connectivity
curl -v http://localhost:8081/health
```

## 🎯 Real-World Examples

### WasteManagement Data
```bash
curl -X POST http://localhost:8081/encrypt \
  -H "Content-Type: application/json" \
  -d '{"uri":"facility/bin123/waste","message":"Bin blood pressure: 120/80, heart rate: 72"}'
```

### Banking Transactions
```bash
curl -X POST http://localhost:8081/encrypt \
  -H "Content-Type: application/json" \
  -d '{"uri":"bank/transfer/12345","message":"Transfer $1500 from checking to savings account"}'
```

### Personal Messages
```bash
curl -X POST http://localhost:8081/encrypt \
  -H "Content-Type: application/json" \
  -d '{"uri":"personal/notes","message":"Meeting with Dr. Smith tomorrow at 3 PM"}'
```

## 📊 What the Metrics Mean

- **executionTime**: How fast the operation completed (microseconds)
- **memoryUsage**: How much memory was used (bytes)
- **cpuPercentage**: CPU usage during operation (%)
- **powerUsage**: Estimated power consumption (watts)
- **energyConsumptionJoules**: Total energy used (joules)

## 🎉 Success!

If you can encrypt and decrypt messages successfully, your HIBE API is working perfectly!

### Next Steps:
- 📖 Read the full [README.md](README.md) for detailed documentation
- 🔧 Check [API_DOCUMENTATION.md](API_DOCUMENTATION.md) for advanced usage
- 🐳 Explore Docker options for production deployment

## 💡 Pro Tips

1. **Use different URIs** for different types of data - they create separate encryption contexts
2. **Monitor the metrics** to understand performance characteristics
3. **Use HTTPS** in production environments
4. **Keep the encrypted data** safe - it can only be decrypted with the correct URI

**🎊 You're now ready to use HIBE API for secure encryption!**