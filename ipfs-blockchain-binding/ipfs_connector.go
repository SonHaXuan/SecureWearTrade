package binding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// IPFSConnector manages IPFS interactions for healthcare data storage
type IPFSConnector struct {
	nodeURL     string
	httpClient  *http.Client
	hashCache   *HashCache
	rateLimiter *IPFSRateLimiter
	mu          sync.RWMutex
}

// HealthcareData represents healthcare data structure for IPFS storage
type HealthcareData struct {
	PatientID    string                 `json:"patient_id"`
	DoctorWallet string                 `json:"doctor_wallet"`
	Department   string                 `json:"department"`
	DataType     string                 `json:"data_type"`
	AccessLevel  string                 `json:"access_level"`
	Timestamp    time.Time              `json:"timestamp"`
	EncryptedData []byte                `json:"encrypted_data"`
	Metadata     map[string]interface{} `json:"metadata"`
	HIBEKeyHash  string                 `json:"hibe_key_hash"`
}

// IPFSResponse represents response from IPFS node
type IPFSResponse struct {
	Hash string `json:"Hash"`
	Name string `json:"Name"`
	Size string `json:"Size"`
}

// HashCache provides caching for IPFS hash operations
type HashCache struct {
	cache   map[string]*CacheEntry
	maxSize int
	mu      sync.RWMutex
}

// IPFSRateLimiter prevents hash flooding attacks
type IPFSRateLimiter struct {
	requests map[string]*RequestTracker
	mu       sync.RWMutex
}

// RequestTracker tracks requests per client
type RequestTracker struct {
	count     int
	resetTime time.Time
	blocked   bool
}

// NewIPFSConnector creates a new IPFS connector
func NewIPFSConnector(nodeURL string) *IPFSConnector {
	return &IPFSConnector{
		nodeURL: nodeURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		hashCache:   NewHashCache(1000),
		rateLimiter: NewIPFSRateLimiter(),
	}
}

// StoreHealthcareData stores encrypted healthcare data on IPFS
func (ic *IPFSConnector) StoreHealthcareData(data *HealthcareData) (string, error) {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	
	// Serialize healthcare data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to serialize healthcare data: %v", err)
	}
	
	// Check rate limiting
	clientID := data.DoctorWallet
	if !ic.rateLimiter.AllowRequest(clientID) {
		return "", fmt.Errorf("rate limit exceeded for client %s", clientID)
	}
	
	// Create multipart form data for IPFS
	body := &bytes.Buffer{}
	boundary := "----WebKitFormBoundary7MA4YWxkTrZu0gW"
	
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Disposition: form-data; name=\"file\"; filename=\"healthcare_data.json\"\r\n")
	body.WriteString("Content-Type: application/json\r\n\r\n")
	body.Write(jsonData)
	body.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	
	// Make request to IPFS
	req, err := http.NewRequest("POST", ic.nodeURL+"/api/v0/add", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))
	
	resp, err := ic.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to store data on IPFS: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("IPFS request failed with status %d", resp.StatusCode)
	}
	
	// Parse IPFS response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read IPFS response: %v", err)
	}
	
	var ipfsResp IPFSResponse
	err = json.Unmarshal(respBody, &ipfsResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse IPFS response: %v", err)
	}
	
	// Cache the hash
	ic.hashCache.Put(ipfsResp.Hash, &CacheEntry{
		Data:       data,
		StoredTime: time.Now(),
	})
	
	return ipfsResp.Hash, nil
}

// RetrieveHealthcareData retrieves healthcare data from IPFS
func (ic *IPFSConnector) RetrieveHealthcareData(hash string, clientID string) (*HealthcareData, error) {
	// Check rate limiting
	if !ic.rateLimiter.AllowRequest(clientID) {
		return nil, fmt.Errorf("rate limit exceeded for client %s", clientID)
	}
	
	// Check cache first
	if cached, found := ic.hashCache.Get(hash); found {
		if data, ok := cached.Data.(*HealthcareData); ok {
			return data, nil
		}
	}
	
	// Retrieve from IPFS
	req, err := http.NewRequest("POST", ic.nodeURL+"/api/v0/cat?arg="+hash, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	resp, err := ic.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data from IPFS: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPFS request failed with status %d", resp.StatusCode)
	}
	
	// Parse healthcare data
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read IPFS response: %v", err)
	}
	
	var data HealthcareData
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse healthcare data: %v", err)
	}
	
	// Cache the result
	ic.hashCache.Put(hash, &CacheEntry{
		Data:       &data,
		StoredTime: time.Now(),
	})
	
	return &data, nil
}

// ValidateHash validates IPFS hash format and existence
func (ic *IPFSConnector) ValidateHash(hash string) (bool, error) {
	// Basic IPFS hash format validation
	if len(hash) != 46 || !strings.HasPrefix(hash, "Qm") {
		return false, fmt.Errorf("invalid IPFS hash format")
	}
	
	// Check if hash exists on IPFS
	req, err := http.NewRequest("POST", ic.nodeURL+"/api/v0/object/stat?arg="+hash, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create validation request: %v", err)
	}
	
	resp, err := ic.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to validate hash on IPFS: %v", err)
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK, nil
}

// PinHash pins a hash to prevent garbage collection
func (ic *IPFSConnector) PinHash(hash string) error {
	req, err := http.NewRequest("POST", ic.nodeURL+"/api/v0/pin/add?arg="+hash, nil)
	if err != nil {
		return fmt.Errorf("failed to create pin request: %v", err)
	}
	
	resp, err := ic.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to pin hash on IPFS: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pin request failed with status %d", resp.StatusCode)
	}
	
	return nil
}

// UnpinHash unpins a hash to allow garbage collection
func (ic *IPFSConnector) UnpinHash(hash string) error {
	req, err := http.NewRequest("POST", ic.nodeURL+"/api/v0/pin/rm?arg="+hash, nil)
	if err != nil {
		return fmt.Errorf("failed to create unpin request: %v", err)
	}
	
	resp, err := ic.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to unpin hash on IPFS: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unpin request failed with status %d", resp.StatusCode)
	}
	
	return nil
}

// GetNodeInfo retrieves IPFS node information
func (ic *IPFSConnector) GetNodeInfo() (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", ic.nodeURL+"/api/v0/id", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create node info request: %v", err)
	}
	
	resp, err := ic.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get node info: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("node info request failed with status %d", resp.StatusCode)
	}
	
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read node info response: %v", err)
	}
	
	var nodeInfo map[string]interface{}
	err = json.Unmarshal(respBody, &nodeInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node info: %v", err)
	}
	
	return nodeInfo, nil
}

// Hash Cache Implementation
func NewHashCache(maxSize int) *HashCache {
	return &HashCache{
		cache:   make(map[string]*CacheEntry),
		maxSize: maxSize,
	}
}

type CacheEntry struct {
	Data       interface{}
	StoredTime time.Time
	AccessCount int64
}

func (hc *HashCache) Get(key string) (*CacheEntry, bool) {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	
	if entry, exists := hc.cache[key]; exists {
		entry.AccessCount++
		return entry, true
	}
	
	return nil, false
}

func (hc *HashCache) Put(key string, entry *CacheEntry) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	
	// Implement LRU eviction if needed
	if len(hc.cache) >= hc.maxSize {
		hc.evictOldest()
	}
	
	hc.cache[key] = entry
}

func (hc *HashCache) evictOldest() {
	oldestKey := ""
	oldestTime := time.Now()
	
	for k, v := range hc.cache {
		if v.StoredTime.Before(oldestTime) {
			oldestTime = v.StoredTime
			oldestKey = k
		}
	}
	
	if oldestKey != "" {
		delete(hc.cache, oldestKey)
	}
}

// Rate Limiter Implementation
func NewIPFSRateLimiter() *IPFSRateLimiter {
	return &IPFSRateLimiter{
		requests: make(map[string]*RequestTracker),
	}
}

func (rl *IPFSRateLimiter) AllowRequest(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	
	tracker, exists := rl.requests[clientID]
	if !exists {
		tracker = &RequestTracker{
			count:     1,
			resetTime: now.Add(time.Minute),
			blocked:   false,
		}
		rl.requests[clientID] = tracker
		return true
	}
	
	// Reset counter if time window passed
	if now.After(tracker.resetTime) {
		tracker.count = 1
		tracker.resetTime = now.Add(time.Minute)
		tracker.blocked = false
		return true
	}
	
	// Check if blocked
	if tracker.blocked {
		return false
	}
	
	// Basic rate limiting (100 requests per minute)
	if tracker.count >= 100 {
		tracker.blocked = true
		return false
	}
	
	tracker.count++
	return true
}

// CleanupExpiredTrackers removes expired request trackers
func (rl *IPFSRateLimiter) CleanupExpiredTrackers() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	for clientID, tracker := range rl.requests {
		if now.After(tracker.resetTime.Add(5 * time.Minute)) {
			delete(rl.requests, clientID)
		}
	}
}

// GetClientStats returns statistics for a specific client
func (rl *IPFSRateLimiter) GetClientStats(clientID string) (int, bool) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	
	if tracker, exists := rl.requests[clientID]; exists {
		return tracker.count, tracker.blocked
	}
	
	return 0, false
}