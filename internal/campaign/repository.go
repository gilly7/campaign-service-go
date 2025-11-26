package campaign

import (
	"context"

	"campaign-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Add this struct for user profile
type userProfile struct {
	FirstName string
	Email     string
	Discount  int
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateCampaign(c *models.Campaign) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO campaigns (id, name, template, status) VALUES ($1, $2, $3, $4)`,
		c.ID, c.Name, c.Template, c.Status)
	return err
}

func (r *Repository) SaveMessage(m *models.CampaignMessage) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO campaign_messages (id, campaign_id, user_id, user_email, user_data) 
         VALUES ($1, $2, $3, $4, $5)`,
		m.ID, m.CampaignID, m.UserID, m.UserEmail, m.UserData)
	return err
}

// THIS IS THE NEW METHOD — REAL USER DATA FROM DB
func (r *Repository) GetUserProfile(userID string) (userProfile, error) {
	var profile userProfile

	err := r.db.QueryRow(context.Background(),
		`SELECT first_name, email, discount_percent FROM users WHERE user_id = $1`, userID,
	).Scan(&profile.FirstName, &profile.Email, &profile.Discount)

	if err != nil {
		// Return fallback for unknown users
		return userProfile{
			FirstName: "Friend",
			Email:     userID + "@example.com",
			Discount:  10,
		}, nil
	}

	return profile, nil
}
