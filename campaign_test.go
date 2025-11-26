package main

import (
	"campaign-service/internal/client"
	"net/http"
	"testing"
	"time"
)

func TestCreateCampaign(t *testing.T) {
	// Give server a moment to start (in real tests you'd wait properly)
	time.Sleep(2 * time.Second)

	resp, err := client.CreateCampaign(
		"http://localhost:8080",
		"Test from Go client",
		"Hi {{.FirstName}}! Welcome aboard!",
		"user1", "user2", "user3",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	t.Log("Campaign created successfully via HTTP client!")
}
