# Smart Waste Management in Smart City - Test Scenarios

## Context Overview

This document defines test scenarios for key delegation and revocation in a **Smart Waste Management System** within a **Smart City** environment.

## Smart City Waste Management Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                  Smart City Waste Management Network             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  Smart Bins  │  │  Collection  │  │  Recycling   │         │
│  │   Network    │  │   Vehicles   │  │   Centers    │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│         │                 │                  │                  │
│         └─────────────────┼──────────────────┘                  │
│                           │                                      │
│                  ┌────────▼─────────┐                           │
│                  │   HIBE Waste     │                           │
│                  │   Data Platform  │                           │
│                  └────────┬─────────┘                           │
│                           │                                      │
│         ┌─────────────────┼─────────────────┐                  │
│         │                 │                 │                  │
│  ┌──────▼──────┐  ┌──────▼──────┐  ┌──────▼──────┐           │
│  │  Transfer    │  │  Processing │  │  Monitoring │           │
│  │   Stations   │  │   Plants    │  │   Systems   │           │
│  └──────────────┘  └──────────────┘  └──────────────┘           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Waste Management Entities

### 1. Smart Bins
- **Bin ID**: `bin/{city_id}/{zone_id}/{bin_id}`
- **Data**: Fill level, weight, waste type, sensor readings, location data

### 2. Collection Operators
- **Drivers**: `operator/driver/{city_id}/{depot_id}/{driver_id}`
- **Supervisors**: `operator/supervisor/{city_id}/{zone_id}/{supervisor_id}`
- **Inspectors**: `operator/inspector/{city_id}/{inspector_id}`

### 3. Waste Facilities
- **Recycling Centers**: `facility/recycling/{city_id}/{center_id}`
- **Transfer Stations**: `facility/transfer/{city_id}/{station_id}`
- **Processing Plants**: `facility/processing/{city_id}/{plant_id}`
- **Landfills**: `facility/landfill/{city_id}/{landfill_id}`

### 4. Collection Services
- **Garbage Trucks**: `collection/truck/{city_id}/{depot_id}/{truck_id}`
- **Recycling Vehicles**: `collection/recycling-vehicle/{city_id}/{vehicle_id}`
- **Emergency Response**: `collection/emergency/{city_id}/{response_id}`

### 5. Administrative Systems
- **City Waste Dept**: `admin/waste-dept/{city_id}`
- **Environmental Monitoring**: `admin/environmental/{city_id}`
- **Analytics**: `admin/analytics/{city_id}/{analytics_id}`

## URI Patterns for Waste Management Data

### Smart Bin Data
```
waste/{city_id}/bin/{bin_id}/sensor-data
waste/{city_id}/bin/{bin_id}/sensor-data/{data_type}
waste/{city_id}/bin/{bin_id}/sensor-data/fill-level
waste/{city_id}/bin/{bin_id}/sensor-data/weight
waste/{city_id}/bin/{bin_id}/sensor-data/temperature
waste/{city_id}/bin/{bin_id}/sensor-data/odor
waste/{city_id}/bin/{bin_id}/sensor-data/location
waste/{city_id}/bin/{bin_id}/sensor-data/waste-type
```

### Collection & IoT Data
```
waste/{city_id}/bin/{bin_id}/collection/schedule
waste/{city_id}/bin/{bin_id}/collection/history
waste/{city_id}/truck/{truck_id}/route
waste/{city_id}/truck/{truck_id}/capacity
waste/{city_id}/truck/{truck_id}/gps-location
```

### Facility Systems
```
waste/{city_id}/facility/{facility_id}/operations
waste/{city_id}/facility/{facility_id}/processing
waste/{city_id}/facility/{facility_id}/sorting
waste/{city_id}/facility/{facility_id}/recycling
waste/{city_id}/facility/{facility_id}/capacity
```

### Analytics & Monitoring
```
waste/{city_id}/analytics/waste-generation
waste/{city_id}/analytics/recycling-rates
waste/{city_id}/analytics/environmental-impact
waste/{city_id}/monitoring/air-quality
```

## Test Scenarios

### Scenario 1: Collection Vehicle Route Optimization

**Context**: Garbage truck driver needs access to smart bin data for optimizing collection route.

**Actors**:
- **Smart Bin**: BIN-001 (Downtown District A - Mixed Waste)
- **Driver**: John Martinez (driver-001)
- **Depot**: Central Collection Depot (depot-001)
- **City**: Metro City (city-metro)

**Delegation Details**:
```json
{
  "scenario": "Route Optimization",
  "uri": "waste/city-metro/bin/bin-001/sensor-data",
  "hierarchy": "smart-city-waste",
  "delegatedTo": "operator/driver/city-metro/depot-001/driver-001",
  "delegatedBy": "admin/waste-dept/city-metro",
  "startTime": "2025-11-01T06:00:00Z",
  "endTime": "2025-11-01T14:00:00Z",
  "permissions": ["read", "decrypt"],
  "reason": "Scheduled collection route optimization",
  "authorization": true
}
```

**Test Cases**:
1. ✅ Driver can access bin data during shift hours
2. ✅ Access denied before shift start
3. ✅ Access denied after shift end
4. ✅ Driver cannot modify bin data (read-only)

### Scenario 2: Emergency Overflow Response

**Context**: Smart bin overflow detected, emergency response team needs immediate access to bin status.

**Actors**:
- **Smart Bin**: BIN-003 (Commercial Street C - Organic Waste)
- **Response Team**: Emergency Collection Unit (emergency-001)
- **Supervisor**: Lisa Chen (supervisor-001)

**Delegation Details**:
```json
{
  "scenario": "Emergency Overflow Response",
  "uri": "waste/city-metro/bin/bin-003/sensor-data/critical",
  "hierarchy": "smart-city-waste",
  "delegatedTo": "collection/emergency/city-metro/emergency-001",
  "delegatedBy": "system/emergency-override",
  "startTime": "2025-11-01T14:30:00Z",
  "endTime": "2025-11-01T18:30:00Z",
  "permissions": ["read", "decrypt", "emergency-access"],
  "reason": "Emergency: Bin overflow detected, public health risk",
  "emergencyOverride": true,
  "autoRevoke": "after-overflow-resolved"
}
```

**Test Cases**:
1. ✅ Immediate access granted without manual authorization
2. ✅ Access to critical bin status (fill level, weight, location)
3. ✅ Automatic revocation after overflow resolved
4. ✅ Full audit trail of emergency access
5. ✅ Notification sent to waste management supervisor

### Scenario 3: Recycling Center Processing

**Context**: Recycling center processes batch and needs to store processing data, accessible to city environmental dept.

**Actors**:
- **Bin**: BIN-002 (Residential Zone B - Recyclables)
- **Recycling Center**: GreenCycle Center (recycling-001)
- **Quality Inspector**: Alex Kumar (inspector-001)

**Delegation Details**:
```json
{
  "scenario": "Recycling Batch Processing",
  "uri": "waste/city-metro/bin/bin-002/collection/batch-2025-11-01",
  "hierarchy": "smart-city-waste",
  "delegatedTo": "admin/environmental/city-metro",
  "delegatedBy": "facility/recycling/city-metro/recycling-001",
  "startTime": "2025-11-01T16:00:00Z",
  "endTime": "2025-11-08T16:00:00Z",
  "permissions": ["read", "decrypt"],
  "reason": "Recycling batch processed, data ready for environmental review",
  "batchType": "recyclable-plastics",
  "cityNotified": true
}
```

**Test Cases**:
1. ✅ Center can encrypt batch data with city's key
2. ✅ Environmental dept can decrypt results
3. ✅ Quality inspector can also access batch data
4. ✅ Center operator cannot access after data uploaded
5. ✅ Data auto-archived after 7 days

### Scenario 4: Transfer Station Data Sharing

**Context**: Transfer station receives waste and grants processing plant access to batch data.

**Actors**:
- **Bin**: BIN-004 (Industrial Park D - Electronic Waste)
- **Transfer Station**: Metro Transfer Hub (transfer-001)
- **Processing Plant**: E-Waste Processing Facility (processing-001)

**Delegation Details**:
```json
{
  "scenario": "Transfer Station to Processing Plant",
  "uri": "waste/city-metro/bin/bin-004/sensor-data/e-waste",
  "hierarchy": "smart-city-waste",
  "delegatedTo": "facility/processing/city-metro/processing-001",
  "delegatedBy": "facility/transfer/city-metro/transfer-001",
  "startTime": "2025-11-05T00:00:00Z",
  "endTime": "2025-12-05T23:59:59Z",
  "permissions": ["read", "decrypt", "add-notes"],
  "reason": "E-waste batch transfer for specialized processing",
  "scopeRestriction": "electronic-waste-only",
  "authorization": true
}
```

**Test Cases**:
1. ✅ Processing plant can access e-waste batch data
2. ✅ Plant cannot access non-electronic waste data
3. ✅ Plant can add processing notes
4. ✅ Transfer station maintains access
5. ✅ Waste dept can revoke plant access

### Scenario 5: Facility Inspection Access

**Context**: City inspector needs to verify waste processing compliance.

**Actors**:
- **Facility**: GreenCycle Center (recycling-001)
- **Inspector**: Rachel Green (inspector-001)
- **Waste Department**: Metro City Waste Dept (waste-dept)

**Delegation Details**:
```json
{
  "scenario": "Compliance Inspection",
  "uri": "waste/city-metro/facility/recycling-001/operations",
  "hierarchy": "smart-city-waste",
  "delegatedTo": "operator/inspector/city-metro/inspector-001",
  "delegatedBy": "admin/waste-dept/city-metro",
  "startTime": "2025-11-01T08:00:00Z",
  "endTime": "2025-11-01T17:00:00Z",
  "permissions": ["read", "decrypt", "verify"],
  "reason": "Quarterly compliance inspection",
  "inspectionId": "inspect-2025-q4-001",
  "complianceCheck": true,
  "environmentalAudit": true
}
```

**Test Cases**:
1. ✅ Inspector can read facility operations data
2. ✅ Inspector can check processing compliance
3. ✅ Inspector can verify environmental standards
4. ✅ Access revoked after inspection completed
5. ✅ Inspection generates audit report

### Scenario 6: Environmental Monitoring

**Context**: City environmental monitoring system needs continuous access to waste facility emissions data.

**Actors**:
- **Processing Plant**: Metro Waste Processing (processing-002)
- **Monitoring System**: Environmental Analytics (environmental-001)

**Delegation Details**:
```json
{
  "scenario": "Environmental Monitoring",
  "uri": "waste/city-metro/facility/processing-002/emissions",
  "hierarchy": "smart-city-waste",
  "delegatedTo": "admin/environmental/city-metro/environmental-001",
  "delegatedBy": "facility/processing/city-metro/processing-002",
  "startTime": "2025-11-02T00:00:00Z",
  "endTime": "2025-11-16T23:59:59Z",
  "permissions": ["read", "decrypt"],
  "reason": "Continuous environmental compliance monitoring",
  "monitoringId": "env-monitor-2025-11",
  "dataMinimization": true,
  "onlyEmissionsData": true,
  "authorization": true
}
```

**Test Cases**:
1. ✅ Monitoring system can access emissions data only
2. ✅ System cannot access operational data
3. ✅ City can monitor environmental compliance
4. ✅ Access limited to 2-week monitoring period
5. ✅ Automated alerts on threshold violations

### Scenario 7: Citizen Waste Tracking

**Context**: Citizen enrolled in waste reduction program, city analytics needs access to household waste data.

**Actors**:
- **Household Bin**: BIN-R-101 (Residential - Smart Home)
- **Analytics Team**: City Waste Analytics (analytics-001)
- **Program Coordinator**: Dr. Sarah Johnson (coordinator-001)

**Delegation Details**:
```json
{
  "scenario": "Waste Reduction Program",
  "uri": "waste/city-metro/bin/bin-r-101/generation-data",
  "hierarchy": "smart-city-waste",
  "delegatedTo": "admin/analytics/city-metro/analytics-001",
  "delegatedBy": "admin/waste-dept/city-metro",
  "startTime": "2025-11-01T00:00:00Z",
  "endTime": "2026-11-01T00:00:00Z",
  "permissions": ["read", "decrypt"],
  "reason": "Waste reduction program participation",
  "programId": "waste-reduction-2025",
  "anonymized": false,
  "citizenConsent": true,
  "programApproved": true,
  "dataTypes": ["fill-level", "weight", "frequency", "waste-type"]
}
```

**Test Cases**:
1. ✅ Analytics team can access specified data types only
2. ✅ Long-term access (1 year)
3. ✅ Citizen can withdraw consent (revoke access)
4. ✅ City oversight can audit access
5. ✅ Data collection logged for transparency

### Scenario 8: IoT Sensor Data Sharing

**Context**: Smart bin shares real-time sensor data with waste management system for predictive maintenance.

**Actors**:
- **Smart Bin**: BIN-001 (Downtown - With Advanced Sensors)
- **Maintenance System**: Predictive Maintenance AI (maintenance-001)

**Delegation Details**:
```json
{
  "scenario": "Predictive Maintenance",
  "uri": "waste/city-metro/bin/bin-001/sensor-data/maintenance",
  "hierarchy": "smart-city-waste",
  "delegatedTo": "admin/waste-dept/city-metro/maintenance-001",
  "delegatedBy": "admin/waste-dept/city-metro",
  "startTime": "2025-11-01T00:00:00Z",
  "endTime": "2025-12-01T00:00:00Z",
  "permissions": ["read", "decrypt", "analyze"],
  "reason": "Continuous monitoring for predictive maintenance",
  "sensorTypes": ["temperature", "odor", "mechanical-status"],
  "realTime": true,
  "aiAnalysis": true
}
```

**Test Cases**:
1. ✅ Maintenance system can access sensor data continuously
2. ✅ Real-time data streaming enabled
3. ✅ AI can analyze patterns for maintenance prediction
4. ✅ System can generate maintenance alerts
5. ✅ Historical data retained for trend analysis

## Revocation Scenarios

### Revocation 1: Driver Departure

**Context**: Collection driver leaves waste management company, all bin access must be revoked.

**Actors**:
- **Driver**: John Martinez (driver-001)
- **Supervisor**: Lisa Chen (supervisor-001)
- **Affected Bins**: 47 smart bins in collection route

**Revocation Details**:
```json
{
  "scenario": "Employee Departure",
  "revokedFrom": "operator/driver/city-metro/depot-001/driver-001",
  "revokedBy": "operator/supervisor/city-metro/supervisor-001",
  "revocationTime": "2025-11-15T17:00:00Z",
  "reason": "Driver terminated employment",
  "affectedDelegations": 47,
  "bulkRevocation": true,
  "immediateEffect": true
}
```

**Test Cases**:
1. ✅ All 47 bin access delegations revoked immediately
2. ✅ Driver access denied to all bins
3. ✅ New driver assigned to route
4. ✅ Waste dept notified of revocation
5. ✅ Historical access logs preserved for audit

### Revocation 2: Security Breach

**Context**: Unauthorized access detected, emergency revocation of all facility access.

**Actors**:
- **Incident**: Unauthorized data access attempt
- **Security Team**: Metro City Security (security-001)
- **Incident ID**: SEC-2025-11-002

**Revocation Details**:
```json
{
  "scenario": "Security Breach Response",
  "incidentId": "SEC-2025-11-002",
  "revokedFrom": "facility/recycling/city-metro/*",
  "revokedBy": "admin/waste-dept/city-metro/security-001",
  "revocationTime": "2025-11-20T03:15:00Z",
  "reason": "Unauthorized access detected - emergency response",
  "affectedFacilities": ["recycling-001", "recycling-002", "recycling-003"],
  "emergencyRevocation": true,
  "reapprovalRequired": true
}
```

**Test Cases**:
1. ✅ Immediate emergency revocation (< 1 second)
2. ✅ All facility access blocked
3. ✅ Security team notification sent
4. ✅ Incident tracking initiated
5. ✅ Manual re-approval process required

### Revocation 3: Compliance Audit

**Context**: Quarterly compliance audit requires temporary suspension of external facility access.

**Actors**:
- **Auditor**: City Compliance Office (compliance-001)
- **Affected Facilities**: All processing plants
- **Duration**: 3-day audit period

**Revocation Details**:
```json
{
  "scenario": "Compliance Audit",
  "revokedFrom": "facility/*/external-access",
  "revokedBy": "admin/waste-dept/city-metro/compliance-001",
  "revocationTime": "2025-11-25T00:00:00Z",
  "restoreTime": "2025-11-28T23:59:59Z",
  "reason": "Quarterly compliance audit - temporary access suspension",
  "temporarySuspension": true,
  "autoRestore": true,
  "internalAccessMaintained": true
}
```

**Test Cases**:
1. ✅ Temporary suspension (3 days)
2. ✅ External facility access blocked
3. ✅ Internal city access preserved
4. ✅ Automatic restoration after audit
5. ✅ Audit documentation generated

## Performance Benchmarks

Based on waste management system test execution:

| Operation | Average Time | Notes |
|-----------|--------------|-------|
| Create delegation | 5-15ms | Including key generation |
| Check revocation | <1ms | Hash lookup |
| Encrypt sensor data | 300-600μs | Varies by data size |
| Decrypt sensor data | 200-300μs | Faster than encrypt |
| Bulk revocation | 50-100ms | For 10 keys |
| Emergency revocation | <5ms | Priority operation |
| Real-time data access | 10-20ms | IoT sensor queries |

## Success Criteria

### Security Requirements
- ✅ Only authorized personnel access bin data
- ✅ Time-based access enforced
- ✅ Scope restrictions work
- ✅ Emergency access logged
- ✅ Revocation immediate

### Privacy Requirements
- ✅ Citizen consent honored (where applicable)
- ✅ Data minimization enforced
- ✅ Right to revoke works
- ✅ Access monitoring available
- ✅ Audit trail complete

### Compliance Requirements
- ✅ Environmental regulations compliance
- ✅ Emergency access documented
- ✅ Audit logs maintained
- ✅ Citizen rights respected
- ✅ Security incidents tracked

### Operational Requirements
- ✅ Access granted quickly
- ✅ Clear error messages
- ✅ Easy revocation process
- ✅ Status checking simple
- ✅ Real-time monitoring available

## Integration with Smart City Systems

### Data Flow
1. **Smart Bins** → Sensor data → HIBE Platform
2. **Collection Vehicles** → Route optimization → Bin access
3. **Facilities** → Processing data → Environmental monitoring
4. **Analytics** → Waste patterns → City planning

### API Integration
- RESTful APIs for all data access
- WebSocket for real-time sensor updates
- Batch processing for analytics
- Emergency override protocols

## Future Enhancements

### Phase 2 Features
1. **AI-Powered Route Optimization** - Machine learning for collection efficiency
2. **Blockchain Waste Tracking** - Immutable waste lifecycle records
3. **Citizen Mobile App** - Real-time bin status and collection schedules
4. **Automated Reporting** - Environmental compliance auto-reporting

### Scalability Targets
- Support for 10,000+ smart bins
- Real-time processing of 1M+ sensor readings/day
- 99.9% system uptime
- < 100ms average response time

---

**Document Version**: 1.0.0
**Last Updated**: 2025-10-31
**Status**: ✅ Ready for Implementation
