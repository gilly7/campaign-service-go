package campaign_test

import (
	"context"
	"testing"

	"campaign-service/internal/message"
	"campaign-service/internal/models"

	"github.com/stretchr/testify/assert"
)

// Simple in-memory queue
type MockQueue struct {
	Messages []models.CampaignMessage
}

func (m *MockQueue) Publish(msg models.CampaignMessage) error {
	m.Messages = append(m.Messages, msg)
	return nil
}

// Mock service using the MockQueue
type testService struct {
	queue *MockQueue
}

func (s *testService) CreateCampaign(ctx context.Context, req models.CreateCampaignRequest) (*models.Campaign, error) {
	campaign := &models.Campaign{
		ID:       "test-123",
		Name:     req.Name,
		Template: req.Template,
		Status:   "active",
	}

	names := []string{"UserA", "UserB", "UserC", "UserD"}

	for i, userID := range req.UserIDs {
		firstName := "UserX"
		if i < len(names) {
			firstName = names[i]
		}

		data := map[string]interface{}{
			"FirstName": firstName,
			"Discount":  20 + (i%3)*10, // ← int
			"OfferURL":  "https://example.com/offer/" + userID,
		}

		msg := models.CampaignMessage{
			ID:         "msg-" + userID,
			CampaignID: campaign.ID,
			UserID:     userID,
			UserEmail:  userID + "@gmail.com",
			UserData:   data,
		}

		if err := s.queue.Publish(msg); err != nil {
			return nil, err
		}
	}

	return campaign, nil
}

func TestCreateCampaign_FullFlow(t *testing.T) {
	queue := &MockQueue{}
	service := &testService{queue: queue}

	req := models.CreateCampaignRequest{
		Name:     "Victory Test",
		Template: "Hi {{.FirstName}}! Get {{.Discount}}% off at {{.OfferURL}}!",
		UserIDs:  []string{"alice", "bob"},
	}

	campaign, err := service.CreateCampaign(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "test-123", campaign.ID)

	assert.Len(t, queue.Messages, 2)

	// Verify first message data
	assert.Equal(t, "UserA", queue.Messages[0].UserData["FirstName"])
	assert.Equal(t, 20, queue.Messages[0].UserData["Discount"]) // ← int, not float64
	assert.Contains(t, queue.Messages[0].UserData["OfferURL"], "alice")

	assert.Equal(t, "UserB", queue.Messages[1].UserData["FirstName"])
	assert.Equal(t, 30, queue.Messages[1].UserData["Discount"]) // ← int

	// Verify template rendering
	templater, err := message.NewTemplater(req.Template)
	assert.NoError(t, err)

	rendered, err := templater.Render(queue.Messages[0].UserData)
	assert.NoError(t, err)
	assert.Equal(t, "Hi UserA! Get 20% off at https://example.com/offer/alice!", rendered)
}
