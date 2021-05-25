package model

import (
	"time"
)

type IPDetails struct {
	CreatedAt time.Time `json:"created_at"`
	UUID      string    `json:"uuid"`
	IPAddress string    `json:"ip_address"`
	ResponseCode *string    `json:"response_code"`
	UpdatedAt     time.Time `json:"updated_at"`
}
