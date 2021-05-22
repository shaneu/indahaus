package model

import (
	"time"
)

// type IPDetails struct {
// 	CreatedAt     time.Time `json:"created_at"`
// 	// the `gqlgen` tag map the ID to uuid in the response
// 	ID            string    `gqlgen"uuid"`
// 	IPAddress     string    `json:"ip_address"`
// 	ResponseCodes string    `json:"response_codes"`
// 	UpdatedAt     time.Time `json:"updated_at"`
// }

type IPDetails struct {
	CreatedAt     time.Time `json:"created_at"`
	UUID          string    `json:"uuid"`
	IPAddress     string    `json:"ip_address"`
	ResponseCodes string    `json:"response_codes"`
	UpdatedAt     time.Time `json:"updated_at"`
}
