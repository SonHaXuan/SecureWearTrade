package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"jedi"
	"math"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/ucbrise/jedi-pairing/lang/go/wkdibe"
)

const TestPatternSize = 20

var TestHierarchy = []byte("testHierarchy")

const quote1 = "Imagination is more important than knowledge. --Albert Einstein"
const quote2 = "Today is your day! / Your mountain is waiting. / So... get on your way! --Theodor Seuss Geisel"

type DecryptRequest struct {
	URI              string `string:"uri" binding:"required"`
	ENCRYPTEDMESSAGE string `string:"encryptedMessage" binding:"required"`
	KEY              string `string:"key" binding:"required"`
}

type EncryptRequest struct {
	URI     string `string:"uri" binding:"required"`
	MESSAGE string `string:"message" binding:"required"`
}

type MeasureUsage struct {
	memory            uint64
	cpuPercentage     float64
	powerUsage        float64  // Power usage in watts
	energyConsumption float64  // Energy consumption in joules
	executionTime     int64    // Execution time in microseconds
}

// PowerConstants represents power consumption coefficients for different operations
type PowerConstants struct {
	cpuFactor    float64 // Watts per % CPU
	memoryFactor float64 // Watts per GB memory
	basePower    float64 // Base power consumption in watts
}

// Default power constants for a typical mobile/IoT device
// These values should be calibrated for specific hardware
var defaultPowerConstants = PowerConstants{
	cpuFactor:    0.05,   // 0.05W per 1% CPU utilization
	memoryFactor: 0.02,   // 0.02W per GB of memory used
	basePower:    0.5,    // 0.5W base power consumption
}

// PowerReport holds data for power consumption analysis
type PowerReport struct {
	OperationName      string    `json:"operationName"`
	Timestamp          time.Time `json:"timestamp"`
	ExecutionTimeMs    int64     `json:"executionTimeMs"`
	PowerUsageWatts    float64   `json:"powerUsageWatts"`
	EnergyJoules       float64   `json:"energyJoules"`
	MemoryUsageKB      uint64    `json:"memoryUsageKB"`
	CPUUtilizationPct  float64   `json:"cpuUtilizationPct"`
	DataSizeBytes      int       `json:"dataSizeBytes,omitempty"`
}

// We'll store historical power reports for analysis
var powerReports = []PowerReport{}

func measureMemoryUsage() MeasureUsage {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Alloc = %v MiB\n", memStats.Alloc/1024/1024)
	fmt.Printf("TotalAlloc = %v MiB\n", memStats.TotalAlloc/1024/1024)
	fmt.Printf("Sys = %v MiB\n", memStats.Sys/1024/1024)
	fmt.Printf("NumGC = %v\n\n", memStats.NumGC)
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println(err)
	}

	// Calculate power usage
	memoryGB := float64(memStats.Alloc) / (1024 * 1024 * 1024)
	cpuPercent := percentages[0]
	powerUsage := estimatePowerConsumption(cpuPercent, memoryGB, defaultPowerConstants)

	return MeasureUsage{
		memory:            memStats.Alloc,
		cpuPercentage:     percentages[0],
		powerUsage:        powerUsage,
		energyConsumption: 0, // Will be calculated after execution time is known
	}
}

// estimatePowerConsumption calculates estimated power consumption based on resource usage
func estimatePowerConsumption(cpuPercent float64, memoryGB float64, constants PowerConstants) float64 {
	cpuPower := cpuPercent * constants.cpuFactor
	memoryPower := memoryGB * constants.memoryFactor

	// Total power is base power plus CPU and memory components
	totalPower := constants.basePower + cpuPower + memoryPower

	// Round to 2 decimal places for readability
	return math.Round(totalPower*100) / 100
}

// calculateEnergyConsumption computes energy (joules) from power (watts) and time (seconds)
func calculateEnergyConsumption(powerWatts float64, timeSeconds float64) float64 {
	energyJoules := powerWatts * timeSeconds
	return math.Round(energyJoules*1000) / 1000 // Round to 3 decimal places
}

// getPowerProfile captures current system-wide resource usage for power profiling
func getPowerProfile() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	v, _ := mem.VirtualMemory()
	cpuPercent, _ := cpu.Percent(0, false)

	memoryGB := float64(memStats.Alloc) / (1024 * 1024 * 1024)
	powerUsage := estimatePowerConsumption(cpuPercent[0], memoryGB, defaultPowerConstants)

	return map[string]interface{}{
		"cpuUtilization":     cpuPercent[0],
		"memoryUsageMB":      memStats.Alloc / (1024 * 1024),
		"totalMemoryMB":      v.Total / (1024 * 1024),
		"estimatedPowerWatts": powerUsage,
	}
}

// generatePowerReport creates a formatted report of power consumption
func generatePowerReport(reports []PowerReport, format string) string {
	if len(reports) == 0 {
		return "No power consumption data available."
	}

	if format == "text" {
		var buffer bytes.Buffer
		buffer.WriteString("============================================\n")
		buffer.WriteString("       POWER CONSUMPTION ANALYSIS REPORT    \n")
		buffer.WriteString("============================================\n\n")
		
		// Summary statistics
		var totalEnergy float64
		var minPower, maxPower float64 = reports[0].PowerUsageWatts, reports[0].PowerUsageWatts
		var avgExecTime int64
		
		for i, report := range reports {
			if i == 0 {
				minPower = report.PowerUsageWatts
				maxPower = report.PowerUsageWatts
			}
			
			totalEnergy += report.EnergyJoules
			avgExecTime += report.ExecutionTimeMs
			
			if report.PowerUsageWatts < minPower {
				minPower = report.PowerUsageWatts
			}
			if report.PowerUsageWatts > maxPower {
				maxPower = report.PowerUsageWatts
			}
			
			// Individual report
			buffer.WriteString(fmt.Sprintf("Operation: %s\n", report.OperationName))
			buffer.WriteString(fmt.Sprintf("Timestamp: %s\n", report.Timestamp.Format(time.RFC3339)))
			buffer.WriteString(fmt.Sprintf("Execution time: %d ms\n", report.ExecutionTimeMs/1000))
			buffer.WriteString(fmt.Sprintf("Power usage: %.2f W\n", report.PowerUsageWatts))
			buffer.WriteString(fmt.Sprintf("Energy consumption: %.3f J\n", report.EnergyJoules))
			buffer.WriteString(fmt.Sprintf("Memory usage: %.2f MB\n", float64(report.MemoryUsageKB)/1024))
			buffer.WriteString(fmt.Sprintf("CPU utilization: %.2f%%\n", report.CPUUtilizationPct))
			if report.DataSizeBytes > 0 {
				buffer.WriteString(fmt.Sprintf("Data size: %d bytes\n", report.DataSizeBytes))
				// Calculate energy efficiency (Joules per byte)
				buffer.WriteString(fmt.Sprintf("Energy efficiency: %.6f J/byte\n", report.EnergyJoules/float64(report.DataSizeBytes)))
			}
			buffer.WriteString("--------------------------------------------\n\n")
		}
		
		// Summary
		buffer.WriteString("============= SUMMARY =============\n")
		buffer.WriteString(fmt.Sprintf("Number of operations: %d\n", len(reports)))
		buffer.WriteString(fmt.Sprintf("Total energy consumed: %.3f J\n", totalEnergy))
		buffer.WriteString(fmt.Sprintf("Power range: %.2f - %.2f W\n", minPower, maxPower))
		buffer.WriteString(fmt.Sprintf("Average execution time: %.2f ms\n", float64(avgExecTime)/float64(len(reports))/1000))
		buffer.WriteString(fmt.Sprintf("Average energy per operation: %.3f J\n", totalEnergy/float64(len(reports))))
		buffer.WriteString("==================================\n")
		
		return buffer.String()
	}
	
	// Default: just return the count of reports
	return fmt.Sprintf("%d power reports available.", len(reports))
}

// recordPowerUsage adds a new power consumption record to the history
func recordPowerUsage(operation string, usage MeasureUsage, dataSize int) {
	report := PowerReport{
		OperationName:      operation,
		Timestamp:          time.Now(),
		ExecutionTimeMs:    usage.executionTime,
		PowerUsageWatts:    usage.powerUsage,
		EnergyJoules:       usage.energyConsumption,
		MemoryUsageKB:      usage.memory / 1024,
		CPUUtilizationPct:  usage.cpuPercentage,
		DataSizeBytes:      dataSize,
	}
	
	// Limit the size of the reports slice to avoid memory bloat
	if len(powerReports) >= 100 {
		// Remove the oldest report (shift left)
		powerReports = powerReports[1:]
	}
	
	powerReports = append(powerReports, report)
}

func main() {
	ctx := context.Background()
	_, store := NewTestKeyStore()
	encoder := jedi.NewDefaultPatternEncoder(TestPatternSize - jedi.MaxTimeLength)

	state := NewTestState()
	now := time.Now()

	r := gin.Default()

	// Add security headers middleware
	r.Use(securityHeaders())

	r.GET("/jedi-private-key", func(c *gin.Context) {
		uri := "a/b/c"

		start := time.Unix(1565119330, 0)
		end := time.Unix(1565219330, 0)

		parent := c.DefaultQuery("parent", "")
		if parent != "" {
			uri := "a/b/c/d"
			delegation, err := jedi.Delegate(ctx, store, encoder, TestHierarchy, uri, start, end, jedi.DecryptPermission|jedi.SignPermission)
			if err != nil {
				fmt.Println(err)
			}

			marshalled := delegation.Marshal()

			c.JSON(200, gin.H{
				"data": marshalled,
			})
			return
		}

		startTime := time.Now()
		delegation, err := jedi.Delegate(ctx, store, encoder, TestHierarchy, uri, start, end, jedi.DecryptPermission|jedi.SignPermission)
		if err != nil {
			fmt.Println(err)
		}
		endTime := time.Now()

		marshalled := delegation.Marshal()

		c.JSON(200, gin.H{
			"time": endTime.Sub(startTime).Microseconds(),
			"data": marshalled,
		})
	})

	r.POST("/encrypt", func(c *gin.Context) {
		measureUsageBefore := measureMemoryUsage()
		fmt.Println(measureUsageBefore)
		var err error

		var encryptRequest EncryptRequest
		c.BindJSON(&encryptRequest)

		message := encryptRequest.MESSAGE
		uri := encryptRequest.URI

		startTime := time.Now()
		var encrypted []byte
		if encrypted, err = state.Encrypt(ctx, TestHierarchy, uri, now, []byte(message)); err != nil {
			fmt.Println(err)
		}
		endTime := time.Now()
		executionTimeMs := endTime.Sub(startTime).Microseconds()

		measureUsage := measureMemoryUsage()
		measureUsage.executionTime = executionTimeMs

		// Calculate energy consumption (execution time in seconds * power usage in watts)
		executionTimeSeconds := float64(executionTimeMs) / 1000000.0
		measureUsage.energyConsumption = calculateEnergyConsumption(measureUsage.powerUsage, executionTimeSeconds)

		// Record power usage with message size for efficiency analysis
		recordPowerUsage("encrypt", measureUsage, len(message))

		var decrypted []byte
		if decrypted, err = state.Decrypt(ctx, TestHierarchy, uri, now, encrypted); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		if !bytes.Equal(decrypted, []byte(message)) {
			fmt.Println("Original and decrypted messages differ")
		}

		c.JSON(200, gin.H{
			"time":                   executionTimeMs,
			"memoryUsage":            measureUsage.memory,
			"cpuPercentage":          measureUsage.cpuPercentage,
			"powerUsageWatts":        measureUsage.powerUsage,
			"energyConsumptionJoules": measureUsage.energyConsumption,
			"data":                   base64.StdEncoding.EncodeToString(encrypted),
		})
	})

	r.POST("/decrypt", func(c *gin.Context) {
		measureUsageBefore := measureMemoryUsage()
		var err error
		var decryptRequest DecryptRequest
		c.BindJSON(&decryptRequest)

		encrypted, err := base64.StdEncoding.DecodeString(decryptRequest.ENCRYPTEDMESSAGE)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		uri := decryptRequest.URI

		startTime := time.Now()
		var decrypted []byte
		if decrypted, err = state.Decrypt(ctx, TestHierarchy, uri, now, encrypted); err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		endTime := time.Now()
		executionTimeMs := endTime.Sub(startTime).Microseconds()

		measureUsage := measureMemoryUsage()
		measureUsage.executionTime = executionTimeMs

		// Calculate energy consumption (execution time in seconds * power usage in watts)
		executionTimeSeconds := float64(executionTimeMs) / 1000000.0
		measureUsage.energyConsumption = calculateEnergyConsumption(measureUsage.powerUsage, executionTimeSeconds)

		// Record power usage with encrypted data size
		recordPowerUsage("decrypt", measureUsage, len(encrypted))

		if !bytes.Equal(decrypted, []byte("test message")) {
			fmt.Println("Original and decrypted messages differ")
		}

		str := string(decrypted)

		c.JSON(200, gin.H{
			"time":                  executionTimeMs,
			"memoryUsage":           measureUsage.memory,
			"cpuPercentage":         measureUsage.cpuPercentage,
			"powerUsageWatts":       measureUsage.powerUsage,
			"energyConsumptionJoules": measureUsage.energyConsumption,
			"data":                  str,
		})
	})

	// Add new endpoint for detailed power analysis
	r.GET("/power-profile", func(c *gin.Context) {
		profile := getPowerProfile()
		c.JSON(200, profile)
	})
	
	// Add endpoint for power consumption reports
	r.GET("/power-report", func(c *gin.Context) {
		format := c.DefaultQuery("format", "json")
		
		if format == "text" {
			report := generatePowerReport(powerReports, "text")
			c.Header("Content-Type", "text/plain")
			c.String(200, report)
			return
		}
		
		// Default to JSON format
		c.JSON(200, gin.H{
			"reportCount": len(powerReports),
			"reports": powerReports,
			"summary": map[string]interface{}{
				"totalOperations": len(powerReports),
				"operationsBreakdown": summarizeOperations(powerReports),
				"averagePowerWatts": calculateAveragePower(powerReports),
				"totalEnergyJoules": calculateTotalEnergy(powerReports),
				"energyEfficiency": calculateEnergyEfficiency(powerReports),
			},
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
	fmt.Println("DONE!")
}

// Add security headers to all responses to improve security posture
func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}

type TestPublicInfo struct {
	params *wkdibe.Params
}

type TestKeyStore struct {
	params *wkdibe.Params
	master *wkdibe.MasterKey
}

func (tpi *TestPublicInfo) ParamsForHierarchy(ctx context.Context, hierarchy []byte) (*wkdibe.Params, error) {
	return tpi.params, nil
}

func testMessageTransfer(state *jedi.ClientState, hierarchy []byte, uri string, timestamp time.Time, message string) {
	var err error
	ctx := context.Background()

	var encrypted []byte
	if encrypted, err = state.Encrypt(ctx, hierarchy, uri, timestamp, []byte(message)); err != nil {
		fmt.Println(err)
	}

	var decrypted []byte
	if decrypted, err = state.Decrypt(ctx, hierarchy, uri, timestamp, encrypted); err != nil {
		fmt.Println(err)
	}

	if !bytes.Equal(decrypted, []byte(message)) {
		fmt.Println("Original and decrypted messages differ")
	}
}

func NewTestKeyStore() (*TestPublicInfo, *TestKeyStore) {
	tks := new(TestKeyStore)
	tks.params, tks.master = wkdibe.Setup(TestPatternSize, true)
	tpi := new(TestPublicInfo)
	tpi.params = tks.params
	return tpi, tks
}

func (tks *TestKeyStore) KeyForPattern(ctx context.Context, hierarchy []byte, pattern jedi.Pattern) (*wkdibe.Params, *wkdibe.SecretKey, error) {
	empty := make(jedi.Pattern, TestPatternSize)
	return tks.params, wkdibe.KeyGen(tks.params, tks.master, empty.ToAttrs()), nil
}

func NewTestState() *jedi.ClientState {
	info, store := NewTestKeyStore()
	encoder := jedi.NewDefaultPatternEncoder(TestPatternSize - jedi.MaxTimeLength)
	return jedi.NewClientState(info, store, encoder, 1<<20)
}

// Helper functions for report generation

func summarizeOperations(reports []PowerReport) map[string]int {
	result := make(map[string]int)
	for _, report := range reports {
		result[report.OperationName]++
	}
	return result
}

func calculateAveragePower(reports []PowerReport) float64 {
	if len(reports) == 0 {
		return 0
	}
	
	var total float64
	for _, report := range reports {
		total += report.PowerUsageWatts
	}
	return math.Round((total/float64(len(reports)))*100) / 100
}

func calculateTotalEnergy(reports []PowerReport) float64 {
	var total float64
	for _, report := range reports {
		total += report.EnergyJoules
	}
	return math.Round(total*1000) / 1000
}

func calculateEnergyEfficiency(reports []PowerReport) map[string]float64 {
	result := make(map[string]float64)
	counts := make(map[string]int)
	
	for _, report := range reports {
		if report.DataSizeBytes > 0 {
			efficiency := report.EnergyJoules / float64(report.DataSizeBytes)
			op := report.OperationName
			result[op] = (result[op]*float64(counts[op]) + efficiency) / float64(counts[op]+1)
			counts[op]++
		}
	}
	
	return result
}
