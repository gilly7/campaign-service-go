package campaign

import (
	"context"
	"log"
	"time"

	"campaign-service/internal/message"
	"campaign-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repo  *Repository
	queue message.Queue
	db    *pgxpool.Pool
}

func NewService(db *pgxpool.Pool, queue message.Queue) *Service {
	return &Service{
		repo:  NewRepository(db),
		queue: queue,
		db:    db,
	}
}

func (s *Service) CreateCampaign(ctx context.Context, req models.CreateCampaignRequest) (*models.Campaign, error) {
	campaign := &models.Campaign{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Template:  req.Template,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateCampaign(campaign); err != nil {
		return nil, err
	}

	for _, userID := range req.UserIDs {
		profile, err := s.repo.GetUserProfile(userID)
		if err != nil {
			log.Printf("Failed to get profile for %s: %v", userID, err)
			continue
		}

		data := map[string]interface{}{
			"FirstName": profile.FirstName,
			"Discount":  profile.Discount,
			"OfferURL":  "https://company.com/offer/" + userID,
		}

		msg := models.CampaignMessage{
			ID:         uuid.New().String(),
			CampaignID: campaign.ID,
			UserID:     userID,
			UserEmail:  profile.Email,
			UserData:   data,
		}

		if err := s.repo.SaveMessage(&msg); err != nil {
			log.Printf("DB save failed: %v", err)
		}
		if err := s.queue.Publish(msg); err != nil {
			log.Printf("Queue failed: %v", err)
		} else {
			log.Printf("Queued message for %s", msg.UserEmail)
		}
	}

	return campaign, nil
}

func (s *Service) GetCampaignByID(id string) (models.Campaign, error) {
	var c models.Campaign
	err := s.db.QueryRow(context.Background(),
		`SELECT id, name, template, status, created_at, scheduled_at 
         FROM campaigns WHERE id = $1`, id).
		Scan(&c.ID, &c.Name, &c.Template, &c.Status, &c.CreatedAt, &c.ScheduledAt)

	if err != nil {
		return c, err
	}
	return c, nil
}

func (s *Service) StartWorker(ctx context.Context) {
	q := s.queue.(*message.RedisQueue)
	q.Consume(ctx, func(msg models.CampaignMessage) error {
		tmpl, err := message.NewTemplater("Hey {{.FirstName}}! Get {{.Discount}}% off – {{.OfferURL}}")
		if err != nil {
			return err
		}
		body, _ := tmpl.Render(msg.UserData)
		log.Printf("Sending to %s: %s", msg.UserEmail, body)
		time.Sleep(100 * time.Millisecond) // simulate send
		log.Printf("Sent to %s", msg.UserEmail)
		return nil
	})
}
