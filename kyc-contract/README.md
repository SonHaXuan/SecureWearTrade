# 🏥 KYC Blockchain Waste Data Trading Platform

A decentralized platform for secure KYC verification and waste data trading using blockchain technology, NFTs, and advanced cryptographic techniques.

## 🌟 Features

### 🔐 Security & Privacy
- **RSA Encryption**: End-to-end encryption for sensitive data
- **HIBE Blockchain Integration**: Advanced cryptographic operations
- **IPFS Storage**: Decentralized metadata storage
- **Owner-Controlled Access**: Granular permission management

### 💎 NFT-Based Data Ownership
- **ERC721 Standard**: Industry-compliant NFT implementation
- **Metadata Storage**: Secure IPFS integration
- **Transfer Mechanism**: Secure data trading with payment
- **Purchase Events**: Transparent transaction tracking

### 🌐 Multi-Network Support
- **Ethereum Sepolia** ✅ (Deployed)
- **BNB Testnet** ✅ (Configured)
- **Fantom Testnet** ✅ (Configured)
- **Celo Testnet** ✅ (Configured)
- **Polygon Testnet** ✅ (Configured)

## 🚀 Quick Start

### 1. Installation
```bash
git clone <repository-url>
cd kyc-contract
npm install
```

### 2. Environment Setup
```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Compile & Deploy
```bash
npx hardhat compile
npx hardhat run scripts/deploy-simple.js --network sepolia
```

## 📖 Documentation

- **[🚀 Deployment Guide](./DEPLOYMENT_GUIDE.md)** - Complete setup and deployment instructions
- **[📚 API Reference](./API_REFERENCE.md)** - Smart contract API documentation
- **[🤖 Claude Guide](./CLAUDE.md)** - Development assistance guide

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   User Data     │    │   IPFS Storage   │    │   KYC Contract  │
│  (Waste/KYC)  │───▶│  (Metadata URI)  │───▶│   (NFT Token)   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│  RSA Keys       │    │  Decentralized   │    │  Blockchain     │
│ (Encryption)    │    │     Storage      │    │  (Ethereum/etc) │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## 🔧 Technology Stack

### Smart Contracts
- **Solidity 0.8.17** - Latest security features
- **OpenZeppelin** - Battle-tested standards
- **Hardhat** - Development framework

### Cryptography
- **RSA-512** - Data encryption (testing)
- **HIBE Blockchain** - Advanced crypto operations
- **IPFS/Pinata** - Decentralized storage

### Development Tools
- **Ethers.js** - Ethereum interaction
- **Chai** - Testing framework
- **TypeScript Support** - Type safety

## 🧪 Testing

```bash
# Run all tests
npx hardhat test

# Run with gas reporting
REPORT_GAS=true npx hardhat test

# Local development
npx hardhat node
npx hardhat test --network localhost
```

### Test Coverage

✅ **RSA Key Generation & Encryption**
✅ **IPFS Metadata Storage**
✅ **NFT Minting & Ownership**
✅ **HIBE Cryptographic Operations**
✅ **Data Transfer & Decryption**

## 📊 Deployed Contracts

### Ethereum Sepolia (Main)
- **Contract**: `0x969c98B11144F58F331a154D002f2Bd53Ee9C2A4`
- **Explorer**: [Etherscan](https://sepolia.etherscan.io/address/0x969c98B11144F58F331a154D002f2Bd53Ee9C2A4)
- **Status**: ✅ Live

### Additional Networks (Ready for Deployment)
- **BNB Testnet** - Configuration complete
- **Fantom Testnet** - Configuration complete
- **Celo Testnet** - Configuration complete
- **Polygon Testnet** - Configuration complete

## 💡 Use Cases

### 🏥 WasteManagement
- Secure bin data sharing
- Waste record verification
- Clinical trial data management
- Operator-bin data exchange

### 🔐 KYC Services
- Identity verification
- Background checks
- Compliance reporting
- Cross-border verification

### 📊 Data Marketplaces
- Personal data trading
- Research data exchange
- Anonymized data sales
- Data monetization

## 🔒 Security Features

### Smart Contract Security
- **Access Control**: Owner-only minting
- **Input Validation**: Comprehensive checks
- **Reentrancy Protection**: Secure transfers
- **Event Logging**: Transparent operations

### Data Security
- **Encryption**: RSA-based encryption
- **Decentralization**: IPFS storage
- **Access Control**: Permission-based access
- **Audit Trail**: Blockchain verification

## 🤝 Contributing

### Development Setup
1. Fork the repository
2. Create feature branch
3. Add tests for new features
4. Ensure all tests pass
5. Submit pull request

### Code Style
- Follow Solidity style guide
- Use TypeScript where possible
- Write comprehensive tests
- Document all functions

## 📝 License

This project is licensed under the ISC License - see the [LICENSE](LICENSE) file for details.

## 🚨 Disclaimer

**This is a proof-of-concept implementation.**

⚠️ **Not for production use** without:
- Security audit
- Legal review
- Compliance checks
- Production-grade encryption (RSA-2048+)

## 📞 Support

### 🐛 Bug Reports
- Open an issue on GitHub
- Provide detailed reproduction steps
- Include error logs and environment details

### 💡 Feature Requests
- Submit feature proposals via GitHub issues
- Include use case and implementation details
- Discuss in community forums

### 📧 General Inquiries
- Create GitHub discussion
- Check existing documentation
- Review deployment guide

---

**Built with ❤️ for secure and transparent data trading**