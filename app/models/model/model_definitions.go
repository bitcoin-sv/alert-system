package model

import "time"

// Misc constants
const (
	gormTypeText               = "text"
	DefaultDatabaseReadTimeout = 20 * time.Second // For all "GET" or "SELECT" methods
)

// All base models
const (
	NameAlertMessage Name = "alert_message" // AlertMessage is the alert message model
	NameEmpty        Name = "empty"         // Empty model (base model without a name set)
	NamePublicKey    Name = "public_key"    // PublicKey is the public key model
)

// All base model table names
const (
	TableAlertMessages = "alert_messages" // TableAlertMessages is the alert message table
	TableEmpty         = "empty"          // TableEmpty is the empty placeholder table
	TablePublicKeys    = "public_keys"    // TablePublicKeys is the public key table
)
