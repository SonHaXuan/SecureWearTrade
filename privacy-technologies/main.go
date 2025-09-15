package main

import (
	"log"
	privacytechnologies "securewear/privacy-technologies"
)

func main() {
	log.Printf("=== SECUREWEAR TRADE: PRIVACY TECHNOLOGIES ANALYSIS SUITE ===")
	log.Printf("Addressing Reviewer #3 Concern: Exploring additional privacy protection technologies")
	log.Printf("Technologies: Differential Privacy, Homomorphic Encryption, Secure Multiparty Computation")
	log.Printf("")
	
	analysisFramework := privacytechnologies.NewPrivacyUtilityAnalysis()
	
	log.Printf("Starting comprehensive privacy-utility analysis for cardiac research scenarios...")
	log.Printf("")
	
	results := analysisFramework.RunComprehensiveAnalysis()
	
	log.Printf("")
	log.Printf("Analysis completed successfully. Generating detailed results...")
	log.Printf("")
	
	results.PrintDetailedResults()
	
	log.Printf("")
	log.Printf("=== PRIVACY TECHNOLOGIES ANALYSIS SUITE COMPLETED ===")
}