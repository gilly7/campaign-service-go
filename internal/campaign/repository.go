package campaign

import (
	"context"

	"campaign-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
