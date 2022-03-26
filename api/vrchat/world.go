package vrchat

import (
	"fmt"
	"net/http"
	"time"

	"github.com/motoki317/vrc-world-tweets-fetcher/utils"
)

type World struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	ImageURL          string `json:"imageUrl"`
	ThumbnailImageURL string `json:"thumbnailImageUrl"`
}

func FetchVRChatWorldInfo(worldID utils.VRChatWorldID) (*World, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	requestURL := fmt.Sprintf("%s/worlds/%s?apiKey=%s", vrchatAPIBasePath, worldID, vrchatAPIKey)
	resp, err := client.Get(requestURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return nil, ErrServer
	} else if resp.StatusCode >= 400 {
		return nil, ErrClient
	}

	var data World
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
