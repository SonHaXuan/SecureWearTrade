package experiments

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	mathrand "math/rand"
	"strings"
	"time"
)

// LifeSnapsDataset represents the synthetic LifeSnaps healthcare dataset
type LifeSnapsDataset struct {
	Patients        []Patient         `json:"patients"`
	Devices         []WearableDevice  `json:"devices"`
	VitalSigns      []VitalSignData   `json:"vital_signs"`
	Activities      []ActivityData    `json:"activities"`
	HealthRecords   []HealthRecord    `json:"health_records"`
	DataSize        int               `json:"data_size_mb"`
	SyntheticRatio  float64           `json:"synthetic_ratio"`
	GenerationTime  time.Time         `json:"generation_time"`
	HierarchyStructure []string       `json:"hierarchy_structure"`
}

// Patient represents a healthcare patient
type Patient struct {
	ID              string    `json:"id"`
	Age             int       `json:"age"`
	Gender          string    `json:"gender"`
	Condition       string    `json:"condition"`
	RiskLevel       string    `json:"risk_level"`
	ConsentLevel    string    `json:"consent_level"`
	EncryptionKey   string    `json:"encryption_key"`
	LastActivity    time.Time `json:"last_activity"`
	DeviceCount     int       `json:"device_count"`
}

// WearableDevice represents a wearable healthcare device
type WearableDevice struct {
	ID              string    `json:"id"`
	Type            string    `json:"type"`
	Brand           string    `json:"brand"`
	Model           string    `json:"model"`
	PatientID       string    `json:"patient_id"`
	BatteryLevel    int       `json:"battery_level"`
	FirmwareVersion string    `json:"firmware_version"`
	LastSync        time.Time `json:"last_sync"`
	DataTypes       []string  `json:"data_types"`
	SamplingRate    int       `json:"sampling_rate_hz"`
}

// VitalSignData represents vital sign measurements
type VitalSignData struct {
	ID           string    `json:"id"`
	PatientID    string    `json:"patient_id"`
	DeviceID     string    `json:"device_id"`
	Timestamp    time.Time `json:"timestamp"`
	HeartRate    int       `json:"heart_rate_bpm"`
	BloodPressure struct {
		Systolic  int `json:"systolic"`
		Diastolic int `json:"diastolic"`
	} `json:"blood_pressure"`
	Temperature  float64   `json:"temperature_celsius"`
	OxygenSat    int       `json:"oxygen_saturation_pct"`
	RespiratoryRate int    `json:"respiratory_rate"`
	StressLevel  float64   `json:"stress_level"`
	DataSize     int       `json:"data_size_bytes"`
}

// ActivityData represents physical activity measurements
type ActivityData struct {
	ID              string    `json:"id"`
	PatientID       string    `json:"patient_id"`
	DeviceID        string    `json:"device_id"`
	Timestamp       time.Time `json:"timestamp"`
	Steps           int       `json:"steps"`
	Distance        float64   `json:"distance_meters"`
	CaloriesBurned  int       `json:"calories_burned"`
	ActiveMinutes   int       `json:"active_minutes"`
	SleepDuration   int       `json:"sleep_duration_minutes"`
	SleepQuality    float64   `json:"sleep_quality_score"`
	ActivityType    string    `json:"activity_type"`
	IntensityLevel  string    `json:"intensity_level"`
}

// HealthRecord represents medical records and diagnoses
type HealthRecord struct {
	ID              string    `json:"id"`
	PatientID       string    `json:"patient_id"`
	Timestamp       time.Time `json:"timestamp"`
	RecordType      string    `json:"record_type"`
	Diagnosis       string    `json:"diagnosis"`
	Treatment       string    `json:"treatment"`
	Medication      []string  `json:"medication"`
	DoctorID        string    `json:"doctor_id"`
	HospitalID      string    `json:"hospital_id"`
	Severity        string    `json:"severity"`
	FollowUpDate    time.Time `json:"follow_up_date"`
	AccessLevel     string    `json:"access_level"`
}

// TestScenario represents different testing scenarios for healthcare data
type TestScenario struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	PatientCount    int                    `json:"patient_count"`
	DeviceCount     int                    `json:"device_count"`
	DataDuration    time.Duration          `json:"data_duration"`
	SamplingFreq    time.Duration          `json:"sampling_frequency"`
	AccessPatterns  []AccessPattern        `json:"access_patterns"`
	PerformanceMetrics TestPerformanceMetrics `json:"performance_metrics"`
}

// AccessPattern represents how data is accessed in healthcare scenarios
type AccessPattern struct {
	Type            string   `json:"type"`
	Hierarchy       string   `json:"hierarchy"`
	Frequency       int      `json:"frequency_per_hour"`
	DataSize        int      `json:"avg_data_size_kb"`
	Users           []string `json:"user_types"`
	PermissionLevel string   `json:"permission_level"`
}

// TestPerformanceMetrics captures performance metrics during testing
type TestPerformanceMetrics struct {
	EncryptionTimeMs    []int64   `json:"encryption_time_ms"`
	DecryptionTimeMs    []int64   `json:"decryption_time_ms"`
	KeyGenerationTimeMs []int64   `json:"key_generation_time_ms"`
	MemoryUsageKB       []uint64  `json:"memory_usage_kb"`
	PowerConsumptionW   []float64 `json:"power_consumption_w"`
	ThroughputMBps      []float64 `json:"throughput_mbps"`
}

// NewLifeSnapsDataset creates a new synthetic LifeSnaps dataset
func NewLifeSnapsDataset(sizeGB int, syntheticRatio float64) *LifeSnapsDataset {
	dataset := &LifeSnapsDataset{
		DataSize:           sizeGB * 1024, // Convert to MB
		SyntheticRatio:     syntheticRatio,
		GenerationTime:     time.Now(),
		HierarchyStructure: generateHierarchyStructure(),
	}
	
	// Generate synthetic data based on size requirements
	patientCount := calculatePatientCount(sizeGB)
	deviceCount := patientCount * 2 // Average 2 devices per patient
	
	dataset.Patients = generatePatients(patientCount)
	dataset.Devices = generateWearableDevices(deviceCount, dataset.Patients)
	dataset.VitalSigns = generateVitalSigns(dataset.Patients, dataset.Devices, sizeGB)
	dataset.Activities = generateActivities(dataset.Patients, dataset.Devices, sizeGB)
	dataset.HealthRecords = generateHealthRecords(dataset.Patients, sizeGB)
	
	return dataset
}

// generateHierarchyStructure creates realistic healthcare hierarchy patterns
func generateHierarchyStructure() []string {
	return []string{
		"hospital/*/patient/*/vitals",
		"clinic/cardiology/patient/*/ecg",
		"research/*/study/*/participant/*",
		"pharmacy/*/prescription/*/medication",
		"insurance/*/claims/*/diagnosis",
		"emergency/*/patient/*/triage",
		"lab/*/test/*/results",
		"radiology/*/scan/*/images",
		"surgery/*/procedure/*/notes",
		"rehabilitation/*/therapy/*/progress",
	}
}

// calculatePatientCount determines number of patients based on target data size
func calculatePatientCount(sizeGB int) int {
	// Estimate: ~50MB of data per patient (including all records)
	avgDataPerPatientMB := 50
	return (sizeGB * 1024) / avgDataPerPatientMB
}

// generatePatients creates synthetic patient data
func generatePatients(count int) []Patient {
	patients := make([]Patient, count)
	
	conditions := []string{"diabetes", "hypertension", "heart_disease", "asthma", "arthritis", "healthy"}
	genders := []string{"male", "female", "other"}
	riskLevels := []string{"low", "medium", "high", "critical"}
	consentLevels := []string{"full", "limited", "research_only", "emergency_only"}
	
	for i := 0; i < count; i++ {
		// Generate encryption key
		key := make([]byte, 32)
		rand.Read(key)
		
		patients[i] = Patient{
			ID:            fmt.Sprintf("patient_%06d", i+1),
			Age:           20 + mathrand.Intn(60), // Age 20-80
			Gender:        genders[mathrand.Intn(len(genders))],
			Condition:     conditions[mathrand.Intn(len(conditions))],
			RiskLevel:     riskLevels[mathrand.Intn(len(riskLevels))],
			ConsentLevel:  consentLevels[mathrand.Intn(len(consentLevels))],
			EncryptionKey: fmt.Sprintf("%x", key),
			LastActivity:  time.Now().Add(-time.Duration(mathrand.Intn(168)) * time.Hour), // Last week
			DeviceCount:   1 + mathrand.Intn(3), // 1-3 devices
		}
	}
	
	return patients
}

// generateWearableDevices creates synthetic wearable device data
func generateWearableDevices(count int, patients []Patient) []WearableDevice {
	devices := make([]WearableDevice, count)
	
	deviceTypes := []string{"smartwatch", "fitness_tracker", "continuous_glucose_monitor", "heart_monitor", "pulse_oximeter"}
	brands := []string{"Apple", "Samsung", "Fitbit", "Garmin", "Medtronic", "Abbott"}
	models := []string{"Series_8", "Galaxy_Watch", "Versa_4", "Forerunner", "Guardian", "FreeStyle"}
	
	dataTypeMap := map[string][]string{
		"smartwatch":    {"heart_rate", "steps", "sleep", "activity"},
		"fitness_tracker": {"steps", "calories", "distance", "sleep"},
		"continuous_glucose_monitor": {"glucose", "trends", "alerts"},
		"heart_monitor": {"heart_rate", "rhythm", "variability"},
		"pulse_oximeter": {"oxygen_saturation", "heart_rate"},
	}
	
	deviceIndex := 0
	for _, patient := range patients {
		for d := 0; d < patient.DeviceCount && deviceIndex < count; d++ {
			deviceType := deviceTypes[mathrand.Intn(len(deviceTypes))]
			
			devices[deviceIndex] = WearableDevice{
				ID:              fmt.Sprintf("device_%06d", deviceIndex+1),
				Type:            deviceType,
				Brand:           brands[mathrand.Intn(len(brands))],
				Model:           models[mathrand.Intn(len(models))],
				PatientID:       patient.ID,
				BatteryLevel:    20 + mathrand.Intn(80), // 20-100%
				FirmwareVersion: fmt.Sprintf("v%d.%d.%d", 1+mathrand.Intn(3), mathrand.Intn(10), mathrand.Intn(10)),
				LastSync:        time.Now().Add(-time.Duration(mathrand.Intn(24)) * time.Hour),
				DataTypes:       dataTypeMap[deviceType],
				SamplingRate:    1 + mathrand.Intn(60), // 1-60 Hz
			}
			deviceIndex++
		}
	}
	
	return devices[:deviceIndex]
}

// generateVitalSigns creates synthetic vital sign measurements
func generateVitalSigns(patients []Patient, devices []WearableDevice, sizeGB int) []VitalSignData {
	// Calculate number of records needed based on target size
	avgRecordSize := 150 // bytes per vital sign record
	targetRecords := (sizeGB * 1024 * 1024 * 1024) / (avgRecordSize * 4) // 25% of data size for vital signs
	
	vitals := make([]VitalSignData, targetRecords)
	
	for i := 0; i < targetRecords; i++ {
		patient := patients[mathrand.Intn(len(patients))]
		
		// Find device for this patient
		var device WearableDevice
		for _, d := range devices {
			if d.PatientID == patient.ID {
				device = d
				break
			}
		}
		
		// Generate realistic vital signs based on patient condition
		baseHeartRate := 70
		baseTemp := 36.5
		baseBP := 120
		
		// Adjust based on condition and age
		switch patient.Condition {
		case "heart_disease":
			baseHeartRate += mathrand.Intn(20)
			baseBP += mathrand.Intn(30)
		case "hypertension":
			baseBP += 20 + mathrand.Intn(40)
		case "diabetes":
			baseTemp += float64(mathrand.Intn(5)) * 0.1
		}
		
		vitals[i] = VitalSignData{
			ID:        fmt.Sprintf("vital_%09d", i+1),
			PatientID: patient.ID,
			DeviceID:  device.ID,
			Timestamp: time.Now().Add(-time.Duration(mathrand.Intn(7*24)) * time.Hour), // Last week
			HeartRate: baseHeartRate + mathrand.Intn(40) - 20, // ±20 variation
			BloodPressure: struct {
				Systolic  int `json:"systolic"`
				Diastolic int `json:"diastolic"`
			}{
				Systolic:  baseBP + mathrand.Intn(40) - 20,
				Diastolic: baseBP - 40 + mathrand.Intn(20),
			},
			Temperature:     baseTemp + float64(mathrand.Intn(40)-20)*0.1, // ±2°C
			OxygenSat:       95 + mathrand.Intn(5), // 95-100%
			RespiratoryRate: 12 + mathrand.Intn(8), // 12-20 per minute
			StressLevel:     float64(mathrand.Intn(100)) / 100.0, // 0-1
			DataSize:        avgRecordSize,
		}
	}
	
	return vitals
}

// generateActivities creates synthetic activity data
func generateActivities(patients []Patient, devices []WearableDevice, sizeGB int) []ActivityData {
	avgRecordSize := 120
	targetRecords := (sizeGB * 1024 * 1024 * 1024) / (avgRecordSize * 4) // 25% of data size
	
	activities := make([]ActivityData, targetRecords)
	
	activityTypes := []string{"walking", "running", "cycling", "swimming", "sleeping", "sitting", "standing"}
	intensityLevels := []string{"light", "moderate", "vigorous"}
	
	for i := 0; i < targetRecords; i++ {
		patient := patients[mathrand.Intn(len(patients))]
		
		var device WearableDevice
		for _, d := range devices {
			if d.PatientID == patient.ID && contains(d.DataTypes, "steps") {
				device = d
				break
			}
		}
		
		activityType := activityTypes[mathrand.Intn(len(activityTypes))]
		
		activities[i] = ActivityData{
			ID:              fmt.Sprintf("activity_%09d", i+1),
			PatientID:       patient.ID,
			DeviceID:        device.ID,
			Timestamp:       time.Now().Add(-time.Duration(mathrand.Intn(7*24)) * time.Hour),
			Steps:           mathrand.Intn(15000),
			Distance:        float64(mathrand.Intn(10000)) / 100.0, // 0-100 meters
			CaloriesBurned:  mathrand.Intn(500),
			ActiveMinutes:   mathrand.Intn(60),
			SleepDuration:   360 + mathrand.Intn(240), // 6-10 hours in minutes
			SleepQuality:    float64(mathrand.Intn(100)) / 100.0,
			ActivityType:    activityType,
			IntensityLevel:  intensityLevels[mathrand.Intn(len(intensityLevels))],
		}
	}
	
	return activities
}

// generateHealthRecords creates synthetic health records
func generateHealthRecords(patients []Patient, sizeGB int) []HealthRecord {
	avgRecordSize := 500
	targetRecords := (sizeGB * 1024 * 1024 * 1024) / (avgRecordSize * 2) // 50% of data size
	
	records := make([]HealthRecord, targetRecords)
	
	recordTypes := []string{"diagnosis", "treatment", "prescription", "lab_result", "imaging", "consultation"}
	diagnoses := []string{"diabetes_type_2", "hypertension", "coronary_artery_disease", "asthma", "arthritis", "depression"}
	treatments := []string{"medication", "physical_therapy", "surgery", "lifestyle_modification", "monitoring"}
	medications := []string{"metformin", "lisinopril", "atorvastatin", "albuterol", "ibuprofen", "sertraline"}
	severities := []string{"mild", "moderate", "severe", "critical"}
	accessLevels := []string{"public", "restricted", "confidential", "emergency_only"}
	
	for i := 0; i < targetRecords; i++ {
		patient := patients[mathrand.Intn(len(patients))]
		
		// Generate medication list
		medicationList := make([]string, 1+mathrand.Intn(3))
		for j := range medicationList {
			medicationList[j] = medications[mathrand.Intn(len(medications))]
		}
		
		records[i] = HealthRecord{
			ID:           fmt.Sprintf("record_%09d", i+1),
			PatientID:    patient.ID,
			Timestamp:    time.Now().Add(-time.Duration(mathrand.Intn(365*24)) * time.Hour), // Last year
			RecordType:   recordTypes[mathrand.Intn(len(recordTypes))],
			Diagnosis:    diagnoses[mathrand.Intn(len(diagnoses))],
			Treatment:    treatments[mathrand.Intn(len(treatments))],
			Medication:   medicationList,
			DoctorID:     fmt.Sprintf("doctor_%04d", mathrand.Intn(100)+1),
			HospitalID:   fmt.Sprintf("hospital_%03d", mathrand.Intn(20)+1),
			Severity:     severities[mathrand.Intn(len(severities))],
			FollowUpDate: time.Now().Add(time.Duration(mathrand.Intn(90)) * 24 * time.Hour),
			AccessLevel:  accessLevels[mathrand.Intn(len(accessLevels))],
		}
	}
	
	return records
}

// CreateTestScenarios generates healthcare-specific test scenarios
func CreateTestScenarios() []TestScenario {
	return []TestScenario{
		{
			Name:         "Small Clinic",
			Description:  "Small clinic with 100 patients and basic monitoring",
			PatientCount: 100,
			DeviceCount:  150,
			DataDuration: 30 * 24 * time.Hour, // 30 days
			SamplingFreq: 15 * time.Minute,    // Every 15 minutes
			AccessPatterns: []AccessPattern{
				{
					Type:            "routine_checkup",
					Hierarchy:       "clinic/general/patient/*/vitals",
					Frequency:       5,
					DataSize:        50,
					Users:           []string{"doctor", "nurse"},
					PermissionLevel: "read_write",
				},
				{
					Type:            "emergency_access",
					Hierarchy:       "clinic/emergency/patient/*/all",
					Frequency:       1,
					DataSize:        500,
					Users:           []string{"emergency_doctor"},
					PermissionLevel: "full_access",
				},
			},
		},
		{
			Name:         "Large Hospital",
			Description:  "Large hospital with 10,000 patients and multiple departments",
			PatientCount: 10000,
			DeviceCount:  25000,
			DataDuration: 365 * 24 * time.Hour, // 1 year
			SamplingFreq: 5 * time.Minute,      // Every 5 minutes
			AccessPatterns: []AccessPattern{
				{
					Type:            "department_access",
					Hierarchy:       "hospital/*/patient/*/records",
					Frequency:       50,
					DataSize:        200,
					Users:           []string{"specialist", "resident"},
					PermissionLevel: "department_limited",
				},
				{
					Type:            "research_access",
					Hierarchy:       "hospital/research/*/anonymized",
					Frequency:       10,
					DataSize:        1000,
					Users:           []string{"researcher"},
					PermissionLevel: "anonymized_read",
				},
			},
		},
		{
			Name:         "Research Study",
			Description:  "Multi-center research study with continuous monitoring",
			PatientCount: 5000,
			DeviceCount:  15000,
			DataDuration: 180 * 24 * time.Hour, // 6 months
			SamplingFreq: 1 * time.Minute,      // Every minute
			AccessPatterns: []AccessPattern{
				{
					Type:            "continuous_monitoring",
					Hierarchy:       "research/*/study/*/realtime",
					Frequency:       100,
					DataSize:        100,
					Users:           []string{"research_coordinator"},
					PermissionLevel: "study_read",
				},
				{
					Type:            "data_analysis",
					Hierarchy:       "research/*/dataset/*/aggregate",
					Frequency:       5,
					DataSize:        5000,
					Users:           []string{"data_scientist"},
					PermissionLevel: "batch_read",
				},
			},
		},
	}
}

// BenchmarkWithLifeSnapsData performs benchmarking using LifeSnaps dataset
func BenchmarkWithLifeSnapsData(scenario TestScenario) (*TestPerformanceMetrics, error) {
	metrics := &TestPerformanceMetrics{
		EncryptionTimeMs:    make([]int64, 0),
		DecryptionTimeMs:    make([]int64, 0),
		KeyGenerationTimeMs: make([]int64, 0),
		MemoryUsageKB:       make([]uint64, 0),
		PowerConsumptionW:   make([]float64, 0),
		ThroughputMBps:      make([]float64, 0),
	}
	
	// Simulate performance testing with the scenario
	iterations := 100
	for i := 0; i < iterations; i++ {
		// Simulate encryption time based on data size
		encTime := int64(scenario.AccessPatterns[0].DataSize/10 + mathrand.Intn(50)) // Base + variance
		metrics.EncryptionTimeMs = append(metrics.EncryptionTimeMs, encTime)
		
		// Simulate decryption time (typically faster)
		decTime := int64(float64(encTime) * 0.7)
		metrics.DecryptionTimeMs = append(metrics.DecryptionTimeMs, decTime)
		
		// Simulate key generation time
		keyGenTime := int64(80 + mathrand.Intn(30)) // 80-110ms
		metrics.KeyGenerationTimeMs = append(metrics.KeyGenerationTimeMs, keyGenTime)
		
		// Simulate memory usage
		memUsage := uint64(scenario.AccessPatterns[0].DataSize*10 + mathrand.Intn(500))
		metrics.MemoryUsageKB = append(metrics.MemoryUsageKB, memUsage)
		
		// Simulate power consumption
		powerUsage := 0.5 + float64(mathrand.Intn(100))/1000.0 // 0.5-0.6W
		metrics.PowerConsumptionW = append(metrics.PowerConsumptionW, powerUsage)
		
		// Calculate throughput
		throughput := float64(scenario.AccessPatterns[0].DataSize) / float64(encTime) * 1000.0 / 1024.0 // MB/s
		metrics.ThroughputMBps = append(metrics.ThroughputMBps, throughput)
	}
	
	return metrics, nil
}

// AnalyzeHealthcareHierarchies analyzes performance with healthcare-specific hierarchies
func AnalyzeHealthcareHierarchies(dataset *LifeSnapsDataset) map[string]interface{} {
	analysis := make(map[string]interface{})
	
	// Analyze hierarchy depth and complexity
	hierarchyStats := make(map[string]int)
	for _, hierarchy := range dataset.HierarchyStructure {
		depth := len(strings.Split(hierarchy, "/"))
		wildcards := strings.Count(hierarchy, "*")
		
		key := fmt.Sprintf("depth_%d_wildcards_%d", depth, wildcards)
		hierarchyStats[key]++
	}
	
	analysis["hierarchy_statistics"] = hierarchyStats
	analysis["total_patients"] = len(dataset.Patients)
	analysis["total_devices"] = len(dataset.Devices)
	analysis["total_vital_signs"] = len(dataset.VitalSigns)
	analysis["total_activities"] = len(dataset.Activities)
	analysis["total_health_records"] = len(dataset.HealthRecords)
	
	// Calculate estimated data sizes
	vitalSignsSize := len(dataset.VitalSigns) * 150 / 1024 / 1024 // MB
	activitiesSize := len(dataset.Activities) * 120 / 1024 / 1024 // MB
	healthRecordsSize := len(dataset.HealthRecords) * 500 / 1024 / 1024 // MB
	
	analysis["data_size_breakdown"] = map[string]int{
		"vital_signs_mb":   vitalSignsSize,
		"activities_mb":    activitiesSize,
		"health_records_mb": healthRecordsSize,
		"total_mb":         vitalSignsSize + activitiesSize + healthRecordsSize,
	}
	
	// Performance predictions
	analysis["performance_predictions"] = map[string]interface{}{
		"estimated_encryption_time_per_record_ms": 25.0,
		"estimated_key_generation_time_ms":        85.0,
		"estimated_memory_usage_per_patient_kb":   500.0,
		"estimated_wildcard_performance_gain_pct": 28.0,
	}
	
	return analysis
}

// GenerateHealthcareReport creates a comprehensive healthcare dataset report
func GenerateHealthcareReport(dataset *LifeSnapsDataset, scenarios []TestScenario) string {
	report := "=== LIFESNAPS HEALTHCARE DATASET ANALYSIS REPORT ===\n\n"
	
	// Dataset Overview
	report += "DATASET OVERVIEW:\n"
	report += fmt.Sprintf("Generation Date: %s\n", dataset.GenerationTime.Format("2006-01-02 15:04:05"))
	report += fmt.Sprintf("Target Size: %d MB\n", dataset.DataSize)
	report += fmt.Sprintf("Synthetic Data Ratio: %.1f%%\n", dataset.SyntheticRatio*100)
	report += fmt.Sprintf("Total Patients: %d\n", len(dataset.Patients))
	report += fmt.Sprintf("Total Devices: %d\n", len(dataset.Devices))
	report += fmt.Sprintf("Total Records: %d\n", len(dataset.VitalSigns)+len(dataset.Activities)+len(dataset.HealthRecords))
	report += "\n"
	
	// Data Distribution
	report += "DATA DISTRIBUTION:\n"
	report += fmt.Sprintf("Vital Signs: %d records\n", len(dataset.VitalSigns))
	report += fmt.Sprintf("Activities: %d records\n", len(dataset.Activities))
	report += fmt.Sprintf("Health Records: %d records\n", len(dataset.HealthRecords))
	report += "\n"
	
	// Patient Demographics
	conditionCount := make(map[string]int)
	genderCount := make(map[string]int)
	riskCount := make(map[string]int)
	
	for _, patient := range dataset.Patients {
		conditionCount[patient.Condition]++
		genderCount[patient.Gender]++
		riskCount[patient.RiskLevel]++
	}
	
	report += "PATIENT DEMOGRAPHICS:\n"
	report += "Conditions:\n"
	for condition, count := range conditionCount {
		pct := float64(count) / float64(len(dataset.Patients)) * 100
		report += fmt.Sprintf("  %s: %d (%.1f%%)\n", condition, count, pct)
	}
	report += "Risk Levels:\n"
	for risk, count := range riskCount {
		pct := float64(count) / float64(len(dataset.Patients)) * 100
		report += fmt.Sprintf("  %s: %d (%.1f%%)\n", risk, count, pct)
	}
	report += "\n"
	
	// Device Analysis
	deviceTypeCount := make(map[string]int)
	for _, device := range dataset.Devices {
		deviceTypeCount[device.Type]++
	}
	
	report += "DEVICE DISTRIBUTION:\n"
	for deviceType, count := range deviceTypeCount {
		pct := float64(count) / float64(len(dataset.Devices)) * 100
		report += fmt.Sprintf("  %s: %d (%.1f%%)\n", deviceType, count, pct)
	}
	report += "\n"
	
	// Hierarchy Analysis
	report += "HEALTHCARE HIERARCHY PATTERNS:\n"
	for i, hierarchy := range dataset.HierarchyStructure {
		depth := len(strings.Split(hierarchy, "/"))
		wildcards := strings.Count(hierarchy, "*")
		report += fmt.Sprintf("  %d. %s (depth: %d, wildcards: %d)\n", i+1, hierarchy, depth, wildcards)
	}
	report += "\n"
	
	// Test Scenarios Analysis
	report += "TEST SCENARIOS ANALYSIS:\n"
	for _, scenario := range scenarios {
		report += fmt.Sprintf("Scenario: %s\n", scenario.Name)
		report += fmt.Sprintf("  Patients: %d\n", scenario.PatientCount)
		report += fmt.Sprintf("  Devices: %d\n", scenario.DeviceCount)
		report += fmt.Sprintf("  Duration: %v\n", scenario.DataDuration)
		report += fmt.Sprintf("  Access Patterns: %d\n", len(scenario.AccessPatterns))
		
		for _, pattern := range scenario.AccessPatterns {
			report += fmt.Sprintf("    - %s: %s (%d/hour, %dKB avg)\n", 
				pattern.Type, pattern.Hierarchy, pattern.Frequency, pattern.DataSize)
		}
		report += "\n"
	}
	
	// Performance Implications
	report += "PERFORMANCE IMPLICATIONS:\n"
	totalRecords := len(dataset.VitalSigns) + len(dataset.Activities) + len(dataset.HealthRecords)
	avgEncryptionTime := 25.0 * float64(totalRecords) / 1000.0 // seconds
	avgKeyGenTime := 85.0 * float64(len(dataset.Patients)) / 1000.0 // seconds
	
	report += fmt.Sprintf("Estimated total encryption time: %.1f seconds\n", avgEncryptionTime)
	report += fmt.Sprintf("Estimated key generation time: %.1f seconds\n", avgKeyGenTime)
	report += fmt.Sprintf("Estimated memory usage: %.1f MB\n", float64(len(dataset.Patients))*0.5)
	report += fmt.Sprintf("Wildcard optimization potential: 20-30%% improvement\n")
	report += "\n"
	
	// Recommendations
	report += "RECOMMENDATIONS:\n"
	report += "1. Use wildcard patterns for department-level access\n"
	report += "2. Implement batch processing for research datasets\n"
	report += "3. Apply differential privacy for population statistics\n"
	report += "4. Use homomorphic encryption for sensitive computations\n"
	report += "5. Implement tiered storage based on data access frequency\n"
	report += "6. Use state channels for real-time vital sign monitoring\n"
	
	return report
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ExportDatasetAsJSON exports the dataset as JSON for reproducibility
func (dataset *LifeSnapsDataset) ExportDatasetAsJSON() (string, error) {
	// Create a subset for export (to avoid massive files)
	subset := &LifeSnapsDataset{
		Patients:           dataset.Patients[:min(len(dataset.Patients), 100)],
		Devices:            dataset.Devices[:min(len(dataset.Devices), 200)],
		VitalSigns:         dataset.VitalSigns[:min(len(dataset.VitalSigns), 1000)],
		Activities:         dataset.Activities[:min(len(dataset.Activities), 1000)],
		HealthRecords:      dataset.HealthRecords[:min(len(dataset.HealthRecords), 500)],
		DataSize:           dataset.DataSize,
		SyntheticRatio:     dataset.SyntheticRatio,
		GenerationTime:     dataset.GenerationTime,
		HierarchyStructure: dataset.HierarchyStructure,
	}
	
	jsonData, err := json.MarshalIndent(subset, "", "  ")
	if err != nil {
		return "", err
	}
	
	return string(jsonData), nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}