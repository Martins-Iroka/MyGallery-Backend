package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/env"
)

type PexelsResponse struct {
	Photos []Photo `json:"photos"`
}

type Photo struct {
	ID           int    `json:"id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	URL          string `json:"url"`
	Photographer string `json:"photographer"`
}

func main() {
	PEXELS_API_URL := "https://api.pexels.com/v1/curated"

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", PEXELS_API_URL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Add("Authorization", env.GetString("PEXELS_API_KEY", ""))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API request failed with status: %s", resp.Status)
		return
	}

	// Read the body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	// Unmarshal JSON into the struct
	var pexelsData PexelsResponse
	err = json.Unmarshal(bodyBytes, &pexelsData)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// Successfully parsed the data!
	log.Printf("Job Success! Fetched %d.", len(pexelsData.Photos))
}
