package ipresult

import (
	"time"
)

// A complete IPResult
type IPResult struct {
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	ID        string    `db:"id" json:"id"`
	IPAddress string    `db:"ip_address" json:"ip_address"`
	// An IP address can have multiple response codes. For now we are simply storing the codes as a
	// comma separated list. This allows the user to see all associated codes for a given ip. Another
	// alternative would be to have a codes table that a result would map to based on a foreign key relationship
	ResponseCode string    `db:"response_code" json:"response_code"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// The subset of fields necessary to construct an IPResult
type NewIPResult struct {
	IPAddress    string `db:"ip_address" json:"ip_address"`
	ResponseCode string `db:"response_code" json:"response_code"`
}

// The subset of fields necessary to update an IPResult
type UpdateIPResult struct {
	ResponseCode string `db:"response_code" json:"response_code"`
}
