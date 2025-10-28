# KYC Contract API Reference

## Smart Contract: KYC

**Contract Address (Sepolia)**: `0x969c98B11144F58F331a154D002f2Bd53Ee9C2A4`

### Contract Overview

The KYC contract implements an ERC721-based system for managing KYC (Know Your Customer) data as NFTs. It allows for secure storage, transfer, and trading of personal data using blockchain technology and IPFS for metadata storage.

### Inheritance

- `ERC721URIStorage` - Provides metadata URI functionality
- `Ownable` - Implements ownership pattern for privileged operations

### State Variables

```solidity
// Token contract reference
IERC20 _token;

// Contract owner address
address public _owner;

// Track sold tokens
mapping(uint256 => bool) public sold;

// Token pricing
mapping(uint256 => uint256) public price;
```

### Events

```solidity
event Purchase(address owner, uint256 price, uint256 id);
```

Emitted when a KYC NFT is purchased/transferred.

## Constructor

```solidity
constructor() ERC721("NFT KYC", "KYC")
```

**Initializes the contract:**
- Sets contract name to "NFT KYC"
- Sets symbol to "KYC"
- Deploys contract owner as initial owner
- Sets token contract reference to owner address

## Public Functions

### mint(uint256 _tokenId, string memory metadataURI)

Mints a new KYC NFT with metadata.

**Access:** `onlyOwner`

**Parameters:**
- `_tokenId` (uint256): Unique identifier for the token
- `metadataURI` (string): IPFS hash containing the KYC metadata

**Returns:** None

**Example:**
```javascript
const tx = await kycContract.mint(1, "QmHash...");
await tx.wait();
```

### transfer(address _to, uint256 _tokenId) external payable

Transfers a KYC NFT to a new owner with payment.

**Access:** External

**Parameters:**
- `_to` (address): Recipient address
- `_tokenId` (uint256): Token ID to transfer

**Payment:** Must include ETH for transfer

**Emits:** `Purchase` event

**Example:**
```javascript
const tx = await kycContract.transfer(
    "0xRecipientAddress...",
    1,
    { value: ethers.utils.parseEther("0.1") }
);
await tx.wait();
```

## View Functions

### ownerOf(uint256 tokenId)

Returns the owner of the specified token.

**Parameters:**
- `tokenId` (uint256): Token ID to query

**Returns:** address - Owner of the token

### tokenURI(uint256 tokenId)

Returns the metadata URI for the specified token.

**Parameters:**
- `tokenId` (uint256): Token ID to query

**Returns:** string - IPFS URI containing metadata

### exists(uint256 tokenId)

Checks if a token exists.

**Parameters:**
- `tokenId` (uint256): Token ID to check

**Returns:** bool - True if token exists

### getApproved(uint256 tokenId)

Returns the approved address for a token.

**Parameters:**
- `tokenId` (uint256): Token ID to query

**Returns:** address - Approved address (or zero address if none)

### isApprovedForAll(address owner, address operator)

Checks if an operator is approved for all tokens of an owner.

**Parameters:**
- `owner` (address): Token owner address
- `operator` (address): Operator address to check

**Returns:** bool - True if operator is approved

## Owner Functions

### owner()

Returns the contract owner address.

**Returns:** address - Contract owner

### transferOwnership(address newOwner)

Transfers contract ownership to a new address.

**Access:** `onlyOwner`

**Parameters:**
- `newOwner` (address): New owner address

### renounceOwnership()

Renounces ownership of the contract.

**Access:** `onlyOwner`

**Warning:** This action is irreversible.

## ERC721 Standard Functions

### balanceOf(address owner)

Returns the number of NFTs owned by an address.

**Parameters:**
- `owner` (address): Address to query

**Returns:** uint256 - Token balance

### approve(address to, uint256 tokenId)

Approves an address to transfer a specific token.

**Parameters:**
- `to` (address): Address to approve
- `tokenId` (uint256): Token ID to approve

### setApprovalForAll(address operator, bool approved)

Sets or unsets approval for an operator to manage all tokens.

**Parameters:**
- `operator` (address): Operator address
- `approved` (bool): Approval status

### safeTransferFrom(address from, address to, uint256 tokenId)

Safely transfers a token from one address to another.

**Parameters:**
- `from` (address): Current owner
- `to` (address): Recipient address
- `tokenId` (uint256): Token ID to transfer

## Usage Examples

### Minting a KYC NFT

```javascript
// Connect to contract
const kycContract = new ethers.Contract(
    contractAddress,
    KYC_ABI,
    ownerSigner
);

// Mint new KYC NFT with IPFS metadata
const tokenId = 1;
const metadataURI = "QmXxx..."; // IPFS hash

const tx = await kycContract.mint(tokenId, metadataURI);
const receipt = await tx.wait();

console.log(`NFT minted: ${receipt.transactionHash}`);
```

### Transferring a KYC NFT

```javascript
// Transfer with payment
const recipientAddress = "0xRecipient...";
const tokenId = 1;
const paymentAmount = ethers.utils.parseEther("0.1");

const tx = await kycContract.transfer(
    recipientAddress,
    tokenId,
    { value: paymentAmount }
);

const receipt = await tx.wait();
console.log(`NFT transferred: ${receipt.transactionHash}`);
```

### Querying Token Information

```javascript
// Get token owner
const owner = await kycContract.ownerOf(1);
console.log(`Token owner: ${owner}`);

// Get token metadata URI
const metadataURI = await kycContract.tokenURI(1);
console.log(`Metadata URI: ${metadataURI}`);

// Get owner's token balance
const balance = await kycContract.balanceOf(owner);
console.log(`Tokens owned: ${balance.toString()}`);
```

## Event Listening

### Listening for Purchase Events

```javascript
kycContract.on("Purchase", (owner, price, tokenId, event) => {
    console.log(`NFT Purchased:`);
    console.log(`  Buyer: ${owner}`);
    console.log(`  Price: ${ethers.utils.formatEther(price)} ETH`);
    console.log(`  Token ID: ${tokenId}`);
    console.log(`  Transaction: ${event.transactionHash}`);
});
```

### Transfer Events

```javascript
kycContract.on("Transfer", (from, to, tokenId) => {
    console.log(`Token transferred from ${from} to ${to}, ID: ${tokenId}`);
});
```

## Error Codes

### Common Error Messages

- `"Error, wrong Token id"` - Token does not exist
- `"You are not the owner of the Nft"` - Only contract owner can mint
- `"ERC721: transfer of token that is not own"` - Transfer from non-owner
- `"ERC721: transfer to zero address"` - Invalid recipient address
- `"insufficient funds"` - Not enough ETH for transaction

## Gas Costs

### Estimated Gas Usage

- `mint`: ~50,000 - 100,000 gas
- `transfer`: ~40,000 - 80,000 gas
- `approve`: ~40,000 - 60,000 gas
- `transferFrom`: ~60,000 - 100,000 gas

### Gas Optimization Tips

1. Batch operations when possible
2. Use low gas price periods
3. Consider using gas tokens for frequent operations
4. Monitor gas prices before transactions

## Security Considerations

### Contract Owner Privileges

The contract owner has significant privileges:
- Can mint new KYC NFTs
- Can set token prices (commented in current version)
- Can control token sales status

### Transfer Validation

The contract includes validation for:
- Token existence
- Ownership verification
- Payment requirements (commented in current version)

### Recommendations

1. Implement multi-signature ownership for production
2. Add time locks for sensitive operations
3. Consider implementing upgradeability patterns
4. Regular security audits recommended