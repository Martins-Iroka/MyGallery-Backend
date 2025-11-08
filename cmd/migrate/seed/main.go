package main

type PexelsResponse struct {
	Photos []Photo `json:"photos"`
}

type Photo struct {
	ID           int    `json:"id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	URL          string `json:"url"`
	Photographer string `json:"photographer"`
	// Add other fields you need, like 'src'
}

func main() {

}
