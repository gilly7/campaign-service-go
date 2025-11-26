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

	for i, userID := range req.UserIDs {
		data := map[string]interface{}{
			"FirstName": "User" + string(rune('A'+i%26)),
			"Discount":  20 + (i%3)*10,
			"OfferURL":  "https://example.com/offer/" + userID,
		}

		msg := models.CampaignMessage{
			ID:         uuid.New().String(),
			CampaignID: campaign.ID,
			UserID:     userID,
			UserEmail:  userID + "@gmail.com",
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
