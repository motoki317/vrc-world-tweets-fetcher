package utils

import (
	"fmt"
)

func BuildTwitterStatusURL(authorUserName, tweetID string) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", authorUserName, tweetID)
}

func BuildVRChatWorldURL(worldID VRChatWorldID) string {
	return fmt.Sprintf("https://vrchat.com/home/world/%s", worldID)
}
