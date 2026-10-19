package repository

import (
	"time"
)

type Website struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	URL             string     `json:"url"`
	IntervalSeconds int        `json:"interval_seconds"`
	TimeoutSeconds  int        `json:"timeout_seconds"`
	IsActive        bool       `json:"is_active"`
	LastStatusCode  *int       `json:"last_status_code"`
	LastCheckedAt   *time.Time `json:"last_checked_at"`
	CreatedAt       time.Time  `json:"created_at"`
	SSLExpiresAt    *time.Time `json:"ssl_expires_at"`
}

type CheckLog struct {
	ID             int       `json:"id"`
	WebsiteID      int       `json:"website_id"`
	StatusCode     *int      `json:"status_code"`
	ResponseTimeMS int       `json:"response_time_ms"`
	IsSuccessful   bool      `json:"is_successful"`
	ErrorMessage   *string   `json:"error_message"`
	CheckedAt      time.Time `json:"checked_at"`
}
