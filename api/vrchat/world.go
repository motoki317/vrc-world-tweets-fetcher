package vrchat

import (
	"fmt"
	"net/http"

	"github.com/motoki317/vrc-world-tweets-fetcher/utils"
)

type World struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	ImageURL          string `json:"imageUrl"`
	ThumbnailImageURL string `json:"thumbnailImageUrl"`
}

func FetchVRChatWorldInfo(worldID utils.VRChatWorldID) (*World, error) {
	requestURL := fmt.Sprintf("%s/worlds/%s?apiKey=%s", vrchatAPIBasePath, worldID, vrchatAPIKey)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
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
