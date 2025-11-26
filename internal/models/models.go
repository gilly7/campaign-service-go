package models

import "time"

type Campaign struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Template    string     `json:"template"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
}

type CampaignMessage struct {
	ID         string                 `json:"id"`
	CampaignID string                 `json:"campaign_id"`
	UserID     string                 `json:"user_id"`
	UserEmail  string                 `json:"user_email"`
	UserData   map[string]interface{} `json:"user_data"`
	Status     string                 `json:"status"`
}

type CreateCampaignRequest struct {
	Name        string   `json:"name"`
	Template    string   `json:"template"`
	UserIDs     []string `json:"user_ids"`
	ScheduledAt *string  `json:"scheduled_at,omitempty"`
}
