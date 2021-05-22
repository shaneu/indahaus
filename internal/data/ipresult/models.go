package ipresult

import (
	"time"
)

// A complete IPResult
type IPResult struct {
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	ID            string    `db:"id" json:"id"`
	IPAddress     string    `db:"ip_address" json:"ip_address"`
	ResponseCodes string    `db:"response_codes" json:"response_codes"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// The subset of fields necessary to construct an IPResult
type NewIPResult struct {
	IPAddress     string `db:"ip_address" json:"ip_address"`
	ResponseCodes string `db:"response_codes" json:"response_codes"`
}

// The subset of fields necessary to update an IPResult
type UpdateIPResult struct {
	ResponseCodes string `db:"response_codes" json:"response_codes"`
}
