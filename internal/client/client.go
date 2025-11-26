package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	// Global client with sane defaults (used by all functions below)
	HTTP = &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
)

// CreateCampaign is a helper you can call from tests or main()
func CreateCampaign(baseURL string, name, template string, userIDs ...string) (*http.Response, error) {
	url := fmt.Sprintf("%s/campaigns", baseURL)
	payload := map[string]any{
		"name":     name,
		"template": template,
		"user_ids": userIDs,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return HTTP.Do(req)
}

// Example shows how to use the CreateCampaign function
func Example() {
	resp, err := CreateCampaign("http://localhost:8080", "Flash Sale", "Hi {{.FirstName}}! Get {{.Discount}}% off!", "alice", "bob", "charlie")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("Campaign created! Status: %d\n", resp.StatusCode)

	// Pretty-print response
	var pretty bytes.Buffer
	json.Indent(&pretty, bodyFrom(resp), "", "  ")
	fmt.Println(pretty.String())
}

func bodyFrom(resp *http.Response) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.Bytes()
}
