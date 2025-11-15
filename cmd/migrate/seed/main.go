package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/db"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/env"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/photo"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/video"
)

type PexelsResponse struct {
	Photos []Photo `json:"photos"`
}

type Photo struct {
	ID           int      `json:"id"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	URL          string   `json:"url"`
	Photographer string   `json:"photographer"`
	PhotoSrc     PhotoSrc `json:"src"`
}

type PhotoSrc struct {
	Original  string `json:"original"`
	Large2x   string `json:"large2x"`
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type PexelsResponse2 struct {
	Video []VideoAPI `json:"videos"`
}

type VideoAPI struct {
	ID          int64                  `json:"id"`
	Video_Image string                 `json:"image"`
	Video_Url   string                 `json:"url"`
	Duration    int                    `json:"duration"`
	Files       []VideoDownloadFileAPI `json:"video_files"`
}

type VideoDownloadFileAPI struct {
	ID         int64  `json:"id"`
	Video_Link string `json:"link"`
	Video_Size int32  `json:"size"`
}

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/mygallery?sslmode=disable")
	conn, err := db.NewPostgreInstance(addr, 3, 3, "15m")
	seed := db.NewSeedInstance(conn)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	PEXELS_API_URL := "https://api.pexels.com/v1/curated?per_page=50"

	client := &http.Client{Timeout: 60 * time.Second}
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

	var photos []photo.PhotoPost
	for _, p := range pexelsData.Photos {
		photo := &photo.PhotoPost{
			ID:           int64(p.ID),
			Photographer: p.Photographer,
			Original:     p.PhotoSrc.Original,
			Large2x:      p.PhotoSrc.Large2x,
			Large:        p.PhotoSrc.Large,
			Medium:       p.PhotoSrc.Medium,
			Small:        p.PhotoSrc.Small,
			Portrait:     p.PhotoSrc.Portrait,
			Landscape:    p.PhotoSrc.Landscape,
			Tiny:         p.PhotoSrc.Tiny,
		}
		photos = append(photos, *photo)
	}

	log.Printf("Photo info lenght %d", len(photos))
	seed.CreatePhotos(photos)

	PEXELS_API_URL = "https://api.pexels.com/videos/popular?per_page=50"

	req, err = http.NewRequest("GET", PEXELS_API_URL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Add("Authorization", env.GetString("PEXELS_API_KEY", ""))

	resp, err = client.Do(req)
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
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	var pexelsData2 PexelsResponse2
	err = json.Unmarshal(bodyBytes, &pexelsData2)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	var videos []video.VideoPostAndDownloadFile
	for _, v := range pexelsData2.Video {
		videoPost := video.VideoPost{
			ID:          v.ID,
			Video_Image: v.Video_Image,
			Video_Url:   v.Video_Url,
			Duration:    v.Duration,
		}
		var videoDFs []video.VideoDownloadFile
		for _, vdf := range v.Files {
			videoDF := &video.VideoDownloadFile{
				ID:            vdf.ID,
				Video_Post_Id: v.ID,
				Video_Link:    vdf.Video_Link,
				Video_Size:    vdf.Video_Size,
			}

			videoDFs = append(videoDFs, *videoDF)
		}
		videoData := video.VideoPostAndDownloadFile{
			VideoPost: videoPost,
			Files:     videoDFs,
		}
		videos = append(videos, videoData)
	}

	seed.CreateVideos(videos)

	log.Printf("Video info %d", len(pexelsData2.Video))
}
