package freshdesk

import (
	"time"
)

type UserResponse struct {
	User User
}
type User struct {
	ID                int                    `json:"id,omitempty"`
	Name              string                 `json:"name,omitempty"`
	Active            bool                   `json:"active"`
	Address           string                 `json:"address,omitempty"`
	CustomerID        string                 `json:"customer_id,omitempty"`
	Deleted           bool                   `json:"deleted,omitempty"`
	Description       string                 `json:"description,omitempty"`
	Email             string                 `json:"email,omitempty"`
	ExternalID        string                 `json:"external_id,omitempty"`
	FacebookProfileID string                 `json:"fb_profile_id,omitempty"`
	JobTitle          string                 `json:"job_title,omitempty"`
	Language          string                 `json:"language,omitempty"`
	Mobile            string                 `json:"mobile,omitempty"`
	Phone             string                 `json:"phone,omitempty"`
	TimeZone          string                 `json:"time_zone,omitempty"`
	TwitterID         string                 `json:"twitter_id,omitempty"`
	CreatedAt         time.Time              `json:"created_at,omitempty"`
	UpdatedAt         time.Time              `json:"updated_at,omitempty"`
	CustomFields      map[string]interface{} `json:"custom_field,omitempty"`
}
