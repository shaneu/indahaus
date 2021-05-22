package model

import (
	"time"
)

type IPDetails struct {
	CreatedAt     time.Time `json:"created_at"`
	ID            string    `gqlgen"uuid"`
	IPAddress     string    `json:"ip_address"`
	ResponseCodes string    `json:"response_codes"`
	UpdatedAt     time.Time `json:"updated_at"`
}
