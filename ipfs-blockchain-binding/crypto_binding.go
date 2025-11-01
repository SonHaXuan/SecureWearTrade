package binding

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
	
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// CryptographicBinding manages IPFS-blockchain key binding operations
type CryptographicBinding struct {
	HIBEKeys      map[string]*HIBEKeyData
	Bindings      map[string]*AccessBinding
	ETHConnector  *EthereumConnector
	IPFSConnector *IPFSConnector
	Cache         *BindingCache
	mu            sync.RWMutex
}

// HIBEKeyData represents hierarchical identity-based encryption key data
type HIBEKeyData struct {
	KeyHex        string    `json:"key_hex"`
	Identity      []string  `json:"identity"`
	Depth         int       `json:"depth"`
	BinID     string    `json:"bin_id"`
	OperatorWallet  string    `json:"operator_wallet"`
	Timestamp     time.Time `json:"timestamp"`
	KeyHash       string    `json:"key_hash"`
}

// AccessBinding represents the cryptographic binding between IPFS and blockchain
type AccessBinding struct {
	BindingHash    string           `json:"binding_hash"`
	HIBEKey        string           `json:"hibe_key"`
	TransactionID  string           `json:"transaction_id"`
	IPFSHash       string           `json:"ipfs_hash"`
	Owner          common.Address   `json:"owner"`
	AccessPolicy   *AccessPolicy    `json:"access_policy"`
	Timestamp      time.Time        `json:"timestamp"`
	IsActive       bool             `json:"is_active"`
	GasFeePaid     *big.Int         `json:"gas_fee_paid"`
}

// AccessPolicy defines blockchain-mediated access control policies
type AccessPolicy struct {
	Owner            common.Address `json:"owner"`
	BinID        string         `json:"bin_id"`
	DataType         string         `json:"data_type"`
	AccessLevel      string         `json:"access_level"`
	Department       string         `json:"department"`
	ExpirationTime   time.Time      `json:"expiration_time"`
	PermittedActions []string       `json:"permitted_actions"`
	GasFeeThreshold  *big.Int       `json:"gas_fee_threshold"`
}

// BindingCache provides LRU caching for frequently accessed bindings
type BindingCache struct {
	cache    map[string]*CacheEntry
	maxSize  int
	hits     int64
	misses   int64
	mu       sync.RWMutex
}

// CacheEntry represents a cached binding entry
type CacheEntry struct {
	Binding   *AccessBinding
	AccessTime time.Time
	HitCount   int64
}

// NewCryptographicBinding creates a new cryptographic binding manager
func NewCryptographicBinding(ethConnector *EthereumConnector, ipfsConnector *IPFSConnector) *CryptographicBinding {
	return &CryptographicBinding{
		HIBEKeys:      make(map[string]*HIBEKeyData),
		Bindings:      make(map[string]*AccessBinding),
		ETHConnector:  ethConnector,
		IPFSConnector: ipfsConnector,
		Cache:         NewBindingCache(1000),
	}
}

// GenerateHIBEKeyForWasteManagement generates HIBE key for waste-management access patterns
func (cb *CryptographicBinding) GenerateHIBEKeyForWasteManagement(binID, operatorWallet, department, dataType, accessLevel string) (*HIBEKeyData, error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	// Algorithm 3: Optimized HIBE Key Generation for WasteManagement Patterns
	identity := []string{"facility", department, "bin", binID, dataType, accessLevel}
	
	// Generate cryptographic key using hierarchical approach
	hibeKey, err := cb.generateHIBEKey(identity)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HIBE key: %v", err)
	}
	
	keyData := &HIBEKeyData{
		KeyHex:       hex.EncodeToString(hibeKey),
		Identity:     identity,
		Depth:        len(identity),
		BinID:    binID,
		OperatorWallet: operatorWallet,
		Timestamp:    time.Now(),
		KeyHash:      cb.calculateKeyHash(hibeKey),
	}
	
	// Store in memory for quick access
	keyID := fmt.Sprintf("%s-%s-%s", binID, operatorWallet, department)
	cb.HIBEKeys[keyID] = keyData
	
	return keyData, nil
}

// CreateCryptographicBinding implements the detailed binding mechanism
func (cb *CryptographicBinding) CreateCryptographicBinding(hibeKeyData *HIBEKeyData, ipfsHash string, gasFeePaid *big.Int) (*AccessBinding, error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	// Step 1: Submit transaction to Ethereum blockchain
	transactionID, err := cb.ETHConnector.SubmitAccessTransaction(hibeKeyData, gasFeePaid)
	if err != nil {
		return nil, fmt.Errorf("failed to submit ETH transaction: %v", err)
	}
	
	// Step 2: Create cryptographic binding
	bindingHash := cb.createBindingHash(hibeKeyData.KeyHex, transactionID)
	
	// Step 3: Create access policy
	accessPolicy := &AccessPolicy{
		Owner:            common.HexToAddress(hibeKeyData.OperatorWallet),
		BinID:        hibeKeyData.BinID,
		DataType:         cb.extractDataType(hibeKeyData.Identity),
		AccessLevel:      cb.extractAccessLevel(hibeKeyData.Identity),
		Department:       cb.extractDepartment(hibeKeyData.Identity),
		ExpirationTime:   time.Now().Add(24 * time.Hour), // 24-hour default
		PermittedActions: []string{"read", "decrypt", "verify"},
		GasFeeThreshold:  gasFeePaid,
	}
	
	// Step 4: Create complete binding
	binding := &AccessBinding{
		BindingHash:   bindingHash,
		HIBEKey:       hibeKeyData.KeyHex,
		TransactionID: transactionID,
		IPFSHash:      ipfsHash,
		Owner:         common.HexToAddress(hibeKeyData.OperatorWallet),
		AccessPolicy:  accessPolicy,
		Timestamp:     time.Now(),
		IsActive:      true,
		GasFeePaid:    gasFeePaid,
	}
	
	// Step 5: Store binding in smart contract
	err = cb.ETHConnector.StoreBinding(binding)
	if err != nil {
		return nil, fmt.Errorf("failed to store binding in smart contract: %v", err)
	}
	
	// Step 6: Cache binding for quick access
	cb.Bindings[bindingHash] = binding
	cb.Cache.Put(bindingHash, binding)
	
	return binding, nil
}

// RealWorldExample demonstrates the complete binding process with concrete values
func (cb *CryptographicBinding) RealWorldExample() error {
	fmt.Println("=== Real-World WasteManagement Data Binding Example ===")
	
	// Scenario: Operator uploads bin sensor data to IPFS
	binID := "12345"
	operatorWallet := "0x742d35Cc6634C0532925a3b8D6Ac6B0ad39CEe5C"
	department := "cardiology"
	dataType := "ecg"
	accessLevel := "realtime"
	ipfsHash := "QmX4e7W8tR9oP2aS6dF3gH5jK8lM9nB1cV4xZ2yA7sE6qT"
	gasFee := big.NewInt(1000000000000000) // 0.001 ETH
	
	fmt.Printf("Step 1: Generate HIBE key for bin %s\n", binID)
	
	// Generate HIBE key
	hibeKeyData, err := cb.GenerateHIBEKeyForWasteManagement(binID, operatorWallet, department, dataType, accessLevel)
	if err != nil {
		return fmt.Errorf("step 1 failed: %v", err)
	}
	
	fmt.Printf("  HIBE Key: %s\n", hibeKeyData.KeyHex[:64]+"...")
	fmt.Printf("  Identity: %v\n", hibeKeyData.Identity)
	fmt.Printf("  Key Hash: %s\n", hibeKeyData.KeyHash)
	
	fmt.Printf("Step 2: Create cryptographic binding\n")
	
	// Create binding
	binding, err := cb.CreateCryptographicBinding(hibeKeyData, ipfsHash, gasFee)
	if err != nil {
		return fmt.Errorf("step 2 failed: %v", err)
	}
	
	fmt.Printf("  Transaction ID: %s\n", binding.TransactionID)
	fmt.Printf("  Binding Hash: %s\n", binding.BindingHash)
	fmt.Printf("  IPFS Hash: %s\n", binding.IPFSHash)
	fmt.Printf("  Gas Fee Paid: %s ETH\n", cb.weiToEth(binding.GasFeePaid))
	
	fmt.Printf("Step 3: Verify binding integrity\n")
	
	// Verify binding
	isValid, err := cb.VerifyBinding(binding.BindingHash)
	if err != nil {
		return fmt.Errorf("step 3 failed: %v", err)
	}
	
	fmt.Printf("  Binding Valid: %t\n", isValid)
	fmt.Printf("  Access Policy Active: %t\n", binding.IsActive)
	
	fmt.Println("✅ Real-world binding example completed successfully!")
	
	return nil
}

// generateHIBEKey generates hierarchical identity-based encryption key
func (cb *CryptographicBinding) generateHIBEKey(identity []string) ([]byte, error) {
	// Simplified HIBE key generation using SHA256-based approach
	hasher := sha256.New()
	
	// Add master secret (in real implementation, this would be more secure)
	hasher.Write([]byte("MASTER_SECRET_KEY_HEALTHCARE_2024"))
	
	// Hash each level of the identity hierarchy
	for i, component := range identity {
		hasher.Write([]byte(fmt.Sprintf("LEVEL_%d_%s", i, component)))
	}
	
	// Generate final key
	key := hasher.Sum(nil)
	
	// Extend key to 64 bytes for stronger security
	hasher.Reset()
	hasher.Write(key)
	hasher.Write([]byte("KEY_EXTENSION"))
	
	extendedKey := append(key, hasher.Sum(nil)...)
	
	return extendedKey, nil
}

// createBindingHash creates the cryptographic binding hash
func (cb *CryptographicBinding) createBindingHash(hibeKeyHex, transactionID string) string {
	// Binding Process as specified:
	// 1. Concatenate: HIBE_key ∥ transaction_ID
	// 2. Hash: SHA256(concatenated_data)
	
	concatenated := hibeKeyHex + transactionID
	hash := sha256.Sum256([]byte(concatenated))
	
	return hex.EncodeToString(hash[:])
}

// calculateKeyHash calculates hash for HIBE key
func (cb *CryptographicBinding) calculateKeyHash(key []byte) string {
	hash := sha256.Sum256(key)
	return hex.EncodeToString(hash[:])
}

// VerifyBinding verifies the integrity of a cryptographic binding
func (cb *CryptographicBinding) VerifyBinding(bindingHash string) (bool, error) {
	// Check local cache first
	if cachedBinding, found := cb.Cache.Get(bindingHash); found {
		return cb.validateBinding(cachedBinding), nil
	}
	
	// Retrieve from blockchain if not cached
	binding, err := cb.ETHConnector.RetrieveBinding(bindingHash)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve binding from blockchain: %v", err)
	}
	
	if binding == nil {
		return false, nil
	}
	
	// Verify binding integrity
	expectedHash := cb.createBindingHash(binding.HIBEKey, binding.TransactionID)
	if expectedHash != bindingHash {
		return false, nil
	}
	
	// Validate access policy
	return cb.validateBinding(binding), nil
}

// validateBinding validates binding and access policy
func (cb *CryptographicBinding) validateBinding(binding *AccessBinding) bool {
	// Check if binding is active
	if !binding.IsActive {
		return false
	}
	
	// Check expiration
	if binding.AccessPolicy.ExpirationTime.Before(time.Now()) {
		return false
	}
	
	// Check gas fee threshold
	if binding.GasFeePaid.Cmp(binding.AccessPolicy.GasFeeThreshold) < 0 {
		return false
	}
	
	return true
}

// Helper functions for extracting components from identity
func (cb *CryptographicBinding) extractDataType(identity []string) string {
	if len(identity) >= 5 {
		return identity[4] // dataType position
	}
	return "unknown"
}

func (cb *CryptographicBinding) extractAccessLevel(identity []string) string {
	if len(identity) >= 6 {
		return identity[5] // accessLevel position
	}
	return "basic"
}

func (cb *CryptographicBinding) extractDepartment(identity []string) string {
	if len(identity) >= 2 {
		return identity[1] // department position
	}
	return "general"
}

// weiToEth converts Wei to ETH for display
func (cb *CryptographicBinding) weiToEth(wei *big.Int) string {
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
	return eth.Text('f', 6)
}

// Cache implementation
func NewBindingCache(maxSize int) *BindingCache {
	return &BindingCache{
		cache:   make(map[string]*CacheEntry),
		maxSize: maxSize,
	}
}

func (bc *BindingCache) Get(key string) (*AccessBinding, bool) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	
	if entry, exists := bc.cache[key]; exists {
		entry.AccessTime = time.Now()
		entry.HitCount++
		bc.hits++
		return entry.Binding, true
	}
	
	bc.misses++
	return nil, false
}

func (bc *BindingCache) Put(key string, binding *AccessBinding) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	
	// Implement simple LRU eviction if needed
	if len(bc.cache) >= bc.maxSize {
		bc.evictOldest()
	}
	
	bc.cache[key] = &CacheEntry{
		Binding:    binding,
		AccessTime: time.Now(),
		HitCount:   0,
	}
}

func (bc *BindingCache) evictOldest() {
	oldestKey := ""
	oldestTime := time.Now()
	
	for k, v := range bc.cache {
		if v.AccessTime.Before(oldestTime) {
			oldestTime = v.AccessTime
			oldestKey = k
		}
	}
	
	if oldestKey != "" {
		delete(bc.cache, oldestKey)
	}
}

func (bc *BindingCache) GetStats() (int64, int64) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.hits, bc.misses
}