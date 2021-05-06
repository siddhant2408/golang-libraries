// Package errorlog provides error log utilities.
package errorlog

import (
	"log"
)

// Print call log.Printf with an improved error formatting.
func Print(err error) {
	log.Printf("Error: %v\n%+v", err, err)
}

// Fatal call log.Fatalf with an improved error formatting.
func Fatal(err error) {
	log.Fatalf("Fatal error: %v\n%+v", err, err)
}
