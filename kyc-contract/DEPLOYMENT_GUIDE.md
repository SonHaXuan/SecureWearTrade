# KYC Smart Contract Deployment Guide

## Overview

This project implements a blockchain-based KYC (Know Your Customer) and medical data trading system using NFTs, IPFS storage, and advanced cryptographic techniques. The smart contract allows users to mint KYC data as NFTs and trade them securely.

## 🏗️ Architecture

### Smart Contract Features
- **ERC721-based NFTs** for KYC data ownership
- **IPFS Integration** for metadata storage
- **Owner-controlled minting** and transfer mechanisms
- **Purchase functionality** for data trading
- **Cryptographic security** with RSA and JEDI integration

### Technology Stack
- **Solidity 0.8.17** for smart contract development
- **Hardhat** for development and deployment
- **OpenZeppelin** for secure contract standards
- **IPFS/Pinata** for decentralized metadata storage
- **RSA Encryption** for data security
- **JEDI Blockchain** for advanced cryptographic operations

## 🚀 Quick Start

### Prerequisites
- Node.js (v14, v16, or v18 recommended)
- npm or yarn
- Wallet with private key
- Testnet tokens for deployment

### Installation

1. **Clone and Install**
```bash
git clone <repository-url>
cd kyc-contract
npm install
```

2. **Environment Setup**
```bash
cp .env.example .env
```

3. **Configure Environment Variables**
Edit `.env` with your configuration:
```env
# Sepolia Network (Already deployed)
SEPOLIA_URL=https://eth-sepolia.g.alchemy.com/v2/YOUR_ALCHEMY_KEY

# Private Key
PRIVATE_KEY=0xyour_private_key_here

# Optional: Explorer API Keys
ETHERSCAN_API_KEY=your_etherscan_api_key
BSCSCAN_API_KEY=your_bscscan_api_key
POLYGONSCAN_API_KEY=your_polygonscan_api_key
FTMSCAN_API_KEY=your_ftmscan_api_key
CELOSCAN_API_KEY=your_celoscan_api_key
```

## 🌐 Supported Networks

| Network | Chain ID | RPC URL | Explorer | Gas Token |
|---------|----------|---------|----------|-----------|
| Ethereum Sepolia | 11155111 | Alchemy/Infura | Etherscan | ETH |
| BNB Testnet | 97 | Binance | BscScan | BNB |
| Fantom Testnet | 4002 | Fantom | FTMScan | FTM |
| Celo Alfajores | 44787 | Celo | CeloScan | CELO |
| Polygon Mumbai | 80001 | Ankr | PolygonScan | MATIC |

## 💰 Getting Testnet Tokens

### BNB Testnet
- **Faucet**: https://testnet.binance.org/faucet-smart
- **Token**: BNB
- **Explorer**: https://testnet.bscscan.com

### Fantom Testnet
- **Faucet**: https://faucet.fantom.network/
- **Token**: FTM
- **Explorer**: https://testnet.ftmscan.com

### Celo Testnet
- **Faucet**: https://faucet.celo.org/
- **Token**: CELO
- **Explorer**: https://alfajores-blockscout.celo-testnet.org

### Polygon Testnet
- **Faucet**: https://faucet.polygon.technology/
- **Token**: MATIC
- **Explorer**: https://mumbai.polygonscan.com

## 📦 Deployment

### Compile Contracts
```bash
npx hardhat compile
```

### Deploy to Specific Network
```bash
# BNB Testnet
npx hardhat run scripts/deploy-simple.js --network bscTestnet

# Fantom Testnet
npx hardhat run scripts/deploy-simple.js --network phantomTestnet

# Celo Testnet
npx hardhat run scripts/deploy-simple.js --network celoTestnet

# Polygon Testnet
npx hardhat run scripts/deploy-simple.js --network polygonTestnet
```

### Deploy to All Networks
```bash
npx hardhat run scripts/deploy-all.js
```

## 🎯 Contract Functions

### Core Functions

#### `mint(uint256 _tokenId, string memory metadataURI)`
- **Purpose**: Mint new KYC NFT (owner only)
- **Parameters**:
  - `_tokenId`: Unique token identifier
  - `metadataURI`: IPFS hash containing KYC data
- **Access**: Contract Owner only

#### `transfer(address _to, uint256 _tokenId) external payable`
- **Purpose**: Transfer KYC NFT with payment
- **Parameters**:
  - `_to`: Recipient address
  - `_tokenId`: Token to transfer
- **Payment**: Requires ETH for transfer

## 🧪 Testing

### Run Test Suite
```bash
# Basic tests
npx hardhat test

# Tests with gas reporting
REPORT_GAS=true npx hardhat test

# Local testing
npx hardhat node
npx hardhat test --network localhost
```

### Test Coverage
- RSA key generation and encryption
- IPFS metadata storage
- NFT minting and ownership
- JEDI cryptographic operations
- Data transfer and decryption

## 📁 Project Structure

```
kyc-contract/
├── contracts/
│   └── KYC.sol              # Main smart contract
├── scripts/
│   ├── deploy.js            # Single network deployment
│   ├── deploy-simple.js     # Simple deployment script
│   └── deploy-all.js        # Multi-network deployment
├── test/
│   └── test.js              # Comprehensive test suite
├── data/
│   └── test.json            # Sample medical data
├── hardhat.config.js        # Hardhat configuration
├── .env.example             # Environment template
├── CLAUDE.md                # Claude Code guidance
└── DEPLOYMENT_GUIDE.md      # This documentation
```

## 🔧 Development Commands

### Development Environment
```bash
# Start local Hardhat node
npx hardhat node

# Deploy to local network
npx hardhat run scripts/deploy.js

# Clean compiled contracts
npx hardhat clean
```

### Verification
```bash
# Verify contract on Etherscan
npx hardhat verify --network <network> <contract-address>

# Example for Sepolia
npx hardhat verify --network sepolia 0x969c98B11144F58F331a154D002f2Bd53Ee9C2A4
```

## 🚨 Security Considerations

### Smart Contract Security
- Contract owner has privileged access to minting functions
- Transfer validation prevents unauthorized transfers
- Payment mechanism for data trading

### Data Security
- RSA encryption for sensitive data (512-bit for testing)
- IPFS for decentralized storage
- JEDI blockchain integration for advanced cryptography

### Development Security
- Never commit private keys to version control
- Use environment variables for sensitive configuration
- Test thoroughly on testnets before mainnet deployment

## 📊 Gas Optimization

### Current Implementation
- Solidity 0.8.17 for latest security features
- Efficient storage patterns
- Minimal external calls

### Recommendations
- Consider using ERC721A for batch minting
- Implement gas-efficient transfer logic
- Optimize metadata storage

## 🔄 Contract Upgrades

Current implementation is **not upgradeable**. For production:
- Consider OpenZeppelin upgradeable contracts
- Implement proxy patterns for seamless upgrades
- Plan governance mechanisms for contract changes

## 📞 Support

### Common Issues
1. **Insufficient funds**: Ensure wallet has testnet tokens
2. **Network connectivity**: Check RPC endpoints
3. **Compilation errors**: Verify Solidity version compatibility
4. **Gas estimation**: Adjust gas limits if needed

### Troubleshooting
```bash
# Check wallet balance
npx hardhat console --network <network-name>
> const balance = await ethers.provider.getBalance("your_address");
> console.log(ethers.utils.formatEther(balance));

# Check contract deployment
npx hardhat run scripts/deploy-simple.js --network <network-name>
```

## 📝 License

This project is licensed under the ISC License.

---

**Deployed Contracts:**
- **Ethereum Sepolia**: `0x969c98B11144F58F331a154D002f2Bd53Ee9C2A4`
- **Explorer**: https://sepolia.etherscan.io/address/0x969c98B11144F58F331a154D002f2Bd53Ee9C2A4