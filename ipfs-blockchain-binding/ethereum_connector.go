package binding

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"
	
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthereumConnector manages blockchain interactions for cryptographic binding
type EthereumConnector struct {
	client      *ethclient.Client
	privateKey  *ecdsa.PrivateKey
	contractABI abi.ABI
	contractAddress common.Address
	chainID     *big.Int
}

// SCTradeContract represents the smart contract structure for Algorithm 1
const SCTradeContractABI = `[
	{
		"inputs": [
			{"name": "bindingHash", "type": "bytes32"},
			{"name": "owner", "type": "address"},
			{"name": "timestamp", "type": "uint256"},
			{"name": "gasFeePaid", "type": "uint256"}
		],
		"name": "storeBinding",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [{"name": "bindingHash", "type": "bytes32"}],
		"name": "getBinding",
		"outputs": [
			{"name": "bindingHash", "type": "bytes32"},
			{"name": "owner", "type": "address"},
			{"name": "timestamp", "type": "uint256"},
			{"name": "isActive", "type": "bool"},
			{"name": "gasFeePaid", "type": "uint256"}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [{"name": "bindingHash", "type": "bytes32"}],
		"name": "deactivateBinding",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getActiveBindingsCount",
		"outputs": [{"name": "", "type": "uint256"}],
		"stateMutability": "view",
		"type": "function"
	}
]`

// NewEthereumConnector creates a new Ethereum blockchain connector
func NewEthereumConnector(nodeURL, privateKeyHex, contractAddressHex string, chainID *big.Int) (*EthereumConnector, error) {
	// Connect to Ethereum node
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}
	
	// Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	
	// Parse contract ABI
	contractABI, err := abi.JSON(strings.NewReader(SCTradeContractABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %v", err)
	}
	
	return &EthereumConnector{
		client:          client,
		privateKey:      privateKey,
		contractABI:     contractABI,
		contractAddress: common.HexToAddress(contractAddressHex),
		chainID:         chainID,
	}, nil
}

// SubmitAccessTransaction submits healthcare access transaction to blockchain
func (ec *EthereumConnector) SubmitAccessTransaction(hibeKeyData *HIBEKeyData, gasFeePaid *big.Int) (string, error) {
	// Get current nonce
	publicKey := ec.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to cast public key to ECDSA")
	}
	
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := ec.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}
	
	// Get current gas price
	gasPrice, err := ec.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %v", err)
	}
	
	// Create transaction data for healthcare access
	txData := map[string]interface{}{
		"patientID":    hibeKeyData.PatientID,
		"doctorWallet": hibeKeyData.DoctorWallet,
		"keyHash":      hibeKeyData.KeyHash,
		"timestamp":    hibeKeyData.Timestamp.Unix(),
		"dataType":     ec.extractDataType(hibeKeyData.Identity),
		"department":   ec.extractDepartment(hibeKeyData.Identity),
	}
	
	// Create transaction
	auth, err := bind.NewKeyedTransactorWithChainID(ec.privateKey, ec.chainID)
	if err != nil {
		return "", fmt.Errorf("failed to create transactor: %v", err)
	}
	
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = gasFeePaid // Send the gas fee as value
	auth.GasLimit = uint64(300000) // Set appropriate gas limit
	auth.GasPrice = gasPrice
	
	// Submit transaction (simplified - in real implementation would call actual contract)
	tx := types.NewTransaction(
		nonce,
		ec.contractAddress,
		gasFeePaid,
		300000,
		gasPrice,
		[]byte(fmt.Sprintf("HEALTHCARE_ACCESS_%s_%s", hibeKeyData.PatientID, hibeKeyData.DoctorWallet)),
	)
	
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(ec.chainID), ec.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}
	
	err = ec.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}
	
	// Return transaction hash
	return signedTx.Hash().Hex(), nil
}

// StoreBinding stores the cryptographic binding in the smart contract
func (ec *EthereumConnector) StoreBinding(binding *AccessBinding) error {
	// Convert binding hash to bytes32
	bindingHashBytes := common.HexToHash(binding.BindingHash)
	
	// Prepare transaction
	auth, err := bind.NewKeyedTransactorWithChainID(ec.privateKey, ec.chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %v", err)
	}
	
	// Set transaction parameters
	auth.Value = binding.GasFeePaid
	auth.GasLimit = uint64(500000)
	
	gasPrice, err := ec.client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}
	auth.GasPrice = gasPrice
	
	// Pack the function call data
	input, err := ec.contractABI.Pack(
		"storeBinding",
		bindingHashBytes,
		binding.Owner,
		big.NewInt(binding.Timestamp.Unix()),
		binding.GasFeePaid,
	)
	if err != nil {
		return fmt.Errorf("failed to pack function call: %v", err)
	}
	
	// Get nonce
	publicKey := ec.privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := ec.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}
	
	// Create and sign transaction
	tx := types.NewTransaction(
		nonce,
		ec.contractAddress,
		binding.GasFeePaid,
		auth.GasLimit,
		auth.GasPrice,
		input,
	)
	
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(ec.chainID), ec.privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}
	
	// Send transaction
	err = ec.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}
	
	// Wait for confirmation
	return ec.waitForConfirmation(signedTx.Hash())
}

// RetrieveBinding retrieves binding from blockchain
func (ec *EthereumConnector) RetrieveBinding(bindingHash string) (*AccessBinding, error) {
	bindingHashBytes := common.HexToHash(bindingHash)
	
	// Pack the function call
	input, err := ec.contractABI.Pack("getBinding", bindingHashBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %v", err)
	}
	
	// Call the contract
	result, err := ec.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &ec.contractAddress,
		Data: input,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}
	
	// Unpack the result
	var bindingResult struct {
		BindingHash [32]byte
		Owner       common.Address
		Timestamp   *big.Int
		IsActive    bool
		GasFeePaid  *big.Int
	}
	
	err = ec.contractABI.UnpackIntoInterface(&bindingResult, "getBinding", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %v", err)
	}
	
	// Convert to AccessBinding struct
	binding := &AccessBinding{
		BindingHash:   common.BytesToHash(bindingResult.BindingHash[:]).Hex(),
		Owner:         bindingResult.Owner,
		Timestamp:     time.Unix(bindingResult.Timestamp.Int64(), 0),
		IsActive:      bindingResult.IsActive,
		GasFeePaid:    bindingResult.GasFeePaid,
		AccessPolicy:  &AccessPolicy{}, // Would be populated from additional contract calls
	}
	
	return binding, nil
}

// DeactivateBinding deactivates a binding in the smart contract
func (ec *EthereumConnector) DeactivateBinding(bindingHash string) error {
	bindingHashBytes := common.HexToHash(bindingHash)
	
	auth, err := bind.NewKeyedTransactorWithChainID(ec.privateKey, ec.chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %v", err)
	}
	
	auth.GasLimit = uint64(200000)
	gasPrice, err := ec.client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}
	auth.GasPrice = gasPrice
	
	// Pack function call
	input, err := ec.contractABI.Pack("deactivateBinding", bindingHashBytes)
	if err != nil {
		return fmt.Errorf("failed to pack function call: %v", err)
	}
	
	// Get nonce and create transaction
	publicKey := ec.privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := ec.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}
	
	tx := types.NewTransaction(
		nonce,
		ec.contractAddress,
		big.NewInt(0),
		auth.GasLimit,
		auth.GasPrice,
		input,
	)
	
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(ec.chainID), ec.privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}
	
	err = ec.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}
	
	return ec.waitForConfirmation(signedTx.Hash())
}

// GetActiveBindingsCount returns the number of active bindings
func (ec *EthereumConnector) GetActiveBindingsCount() (*big.Int, error) {
	input, err := ec.contractABI.Pack("getActiveBindingsCount")
	if err != nil {
		return nil, fmt.Errorf("failed to pack function call: %v", err)
	}
	
	result, err := ec.client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &ec.contractAddress,
		Data: input,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}
	
	var count *big.Int
	err = ec.contractABI.UnpackIntoInterface(&count, "getActiveBindingsCount", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %v", err)
	}
	
	return count, nil
}

// waitForConfirmation waits for transaction confirmation
func (ec *EthereumConnector) waitForConfirmation(txHash common.Hash) error {
	timeout := time.After(2 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-timeout:
			return fmt.Errorf("transaction confirmation timeout")
		case <-ticker.C:
			receipt, err := ec.client.TransactionReceipt(context.Background(), txHash)
			if err != nil {
				continue // Transaction not yet mined
			}
			
			if receipt.Status == 1 {
				return nil // Success
			} else {
				return fmt.Errorf("transaction failed with status %d", receipt.Status)
			}
		}
	}
}

// Helper functions
func (ec *EthereumConnector) extractDataType(identity []string) string {
	if len(identity) >= 5 {
		return identity[4]
	}
	return "unknown"
}

func (ec *EthereumConnector) extractDepartment(identity []string) string {
	if len(identity) >= 2 {
		return identity[1]
	}
	return "general"
}

// GetTransactionByHash retrieves transaction details
func (ec *EthereumConnector) GetTransactionByHash(txHash string) (*types.Transaction, error) {
	hash := common.HexToHash(txHash)
	tx, _, err := ec.client.TransactionByHash(context.Background(), hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}
	return tx, nil
}

// EstimateGasForBinding estimates gas cost for storing a binding
func (ec *EthereumConnector) EstimateGasForBinding(binding *AccessBinding) (uint64, error) {
	bindingHashBytes := common.HexToHash(binding.BindingHash)
	
	input, err := ec.contractABI.Pack(
		"storeBinding",
		bindingHashBytes,
		binding.Owner,
		big.NewInt(binding.Timestamp.Unix()),
		binding.GasFeePaid,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to pack function call: %v", err)
	}
	
	publicKey := ec.privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	
	gasLimit, err := ec.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: fromAddress,
		To:   &ec.contractAddress,
		Data: input,
	})
	
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %v", err)
	}
	
	return gasLimit, nil
}