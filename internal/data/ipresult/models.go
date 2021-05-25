package ipresult

import (
	"time"
)

// A complete IPResult
type IPResult struct {
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	ID        string    `db:"id" json:"id"`
	IPAddress string    `db:"ip_address" json:"ip_address"`
	// Storing response code as a pointer to represent nil when an IP address has zero codes
	ResponseCode *string   `db:"response_code" json:"response_code"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// The subset of fields necessary to construct an IPResult
type NewIPResult struct {
	IPAddress    string  `db:"ip_address" json:"ip_address"`
	ResponseCode *string `db:"response_code" json:"response_code"`
}

// The subset of fields necessary to update an IPResult
type UpdateIPResult struct {
	ResponseCode *string `db:"response_code" json:"response_code"`
}
