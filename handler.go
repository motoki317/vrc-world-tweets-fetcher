package main

import (
	"fmt"
	"log"

	"github.com/Ras96/traq-writer"
	"github.com/gofrs/uuid"

	"github.com/motoki317/vrc-world-tweets-fetcher/api/vrchat"
	"github.com/motoki317/vrc-world-tweets-fetcher/utils"
)

type handlerFunc func(tweetURL string, worldID utils.VRChatWorldID, worldInfo *vrchat.World)

func handlerStdoutLogger(tweetURL string, _ utils.VRChatWorldID, worldInfo *vrchat.World) {
	log.Printf(`Discovered a new world!
Name: %s
Description: %s
Original tweet: %s
`, worldInfo.Name, worldInfo.Description, tweetURL)
}

func handlerTraqWebhookLogger(origin string, webhookID uuid.UUID, secret string) handlerFunc {
	w := traqwriter.NewTraqWebhookWriter(webhookID.String(), secret, origin)

	return func(tweetURL string, worldID utils.VRChatWorldID, worldInfo *vrchat.World) {
		_, err := fmt.Fprintf(w, `### %s

%s

%s
%s`, worldInfo.Name, worldInfo.Description, utils.BuildVRChatWorldURL(worldID), tweetURL)
		if err != nil {
			log.Printf("Encountered an error while posting to traQ (origin: %s): %s", origin, err)
		}
	}
}
