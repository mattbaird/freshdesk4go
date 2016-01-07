package freshdesk

import (
	"time"
)

type CustomerResponse struct {
	Customer Customer
}
type Customer struct {
	ID           int                    `json:"id,omitempty"`
	Name         string                 `json:"name,omitempty"`
	CustomerID   string                 `json:"cust_identifier,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Domains      string                 `json:"domains,omitempty"`
	Email        string                 `json:"email,omitempty"`
	Note         string                 `json:"note,omitempty"`
	SLAPolicyID  int                    `json:"sla_policy_id,omitempty"`
	CreatedAt    time.Time              `json:"created_at,omitempty"`
	UpdatedAt    time.Time              `json:"updated_at,omitempty"`
	CustomFields map[string]interface{} `json:"custom_field,omitempty"`
}
