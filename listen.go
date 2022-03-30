package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/gofrs/uuid"
	"github.com/samber/lo"
	"github.com/sivchari/gotwtr"

	"github.com/motoki317/vrc-world-tweets-fetcher/api/twitter"
	"github.com/motoki317/vrc-world-tweets-fetcher/api/vrchat"
	"github.com/motoki317/vrc-world-tweets-fetcher/db/migrate"
	"github.com/motoki317/vrc-world-tweets-fetcher/db/model"
	"github.com/motoki317/vrc-world-tweets-fetcher/db/repository"
	"github.com/motoki317/vrc-world-tweets-fetcher/db/repository/gorm"
	"github.com/motoki317/vrc-world-tweets-fetcher/utils"
)

type listener struct {
	handlers []handlerFunc
	repo     repository.Repository
}

func (l *listener) processFoundVRChatWorld(tweetURL string, worldID utils.VRChatWorldID) error {
	worldURL := utils.BuildVRChatWorldURL(worldID)
	log.Printf("Found an embedded world in tweet: %s\n", worldURL)

	created, err := l.repo.CreateWorldIfNotExists(&model.World{
		ID:            worldID,
		OriginalTweet: tweetURL,
	})
	if err != nil {
		return fmt.Errorf("encountered an error while accessing world repository: %w", err)
	}
	newlyFound := created
	if !newlyFound {
		log.Println("Known world, skipping process.")
		return nil
	}

	// process
	worldInfo, err := vrchat.FetchVRChatWorldInfo(worldID)
	if err != nil {
		return fmt.Errorf("encountered an error while fetching VRChat world info for world %s: %w\n", worldID, err)
	}
	var wg sync.WaitGroup
	for _, h := range l.handlers {
		wg.Add(1)
		go func(h handlerFunc) {
			h(tweetURL, worldID, worldInfo)
			wg.Done()
		}(h)
	}
	wg.Wait()
	return nil
}

func retrieveAuthorUserName(res *gotwtr.ConnectToStreamResponse) (string, error) {
	for _, user := range res.Includes.Users {
		if user.ID == res.Tweet.AuthorID {
			return user.UserName, nil
		}
	}
	return "", errors.New("failed to retrieve author name")
}

func (l *listener) processTweet(event *gotwtr.ConnectToStreamResponse) {
	tweet := event.Tweet
	// check required fields to process
	if tweet == nil {
		log.Println("Received event with nil tweet field. Perhaps something went wrong?")
		return
	}
	if event.Includes == nil {
		log.Println("Received event with nil includes field. Perhaps something went wrong?")
		return
	}
	if tweet.Entities == nil {
		log.Println("Received event with nil tweet.entities field. Perhaps something went wrong?")
		return
	}

	// retrieve author name
	authorUserName, err := retrieveAuthorUserName(event)
	if err != nil {
		log.Println(err)
		return
	}

	// process
	tweetURL := utils.BuildTwitterStatusURL(authorUserName, tweet.ID)
	log.Printf("New tweet (%s): %s\n", tweetURL, tweet.Text)

	if tweet.Entities == nil {
		log.Println("Received a tweet without the entities field, skipping URL extraction.")
		return
	}

	worldIDs := utils.ExtractVRChatWorldIDs(lo.Map[*gotwtr.TweetURL, string](tweet.Entities.URLs, func(url *gotwtr.TweetURL, _ int) string { return url.ExpandedURL }))
	for _, worldID := range worldIDs {
		if err := l.processFoundVRChatWorld(tweetURL, worldID); err != nil {
			log.Printf("An error occurred while processing found world: %s\n", err)
		}
	}
}

func initHandlers() (handlers []handlerFunc, err error) {
	fullSpec := os.Getenv("HANDLERS")
	if fullSpec == "" {
		log.Println("HANDLERS environment variable is empty or not set. Defaults to stdout handler.")
		return []handlerFunc{handlerStdoutLogger}, nil
	}

	specs := strings.Split(fullSpec, ",")
	for _, spec := range specs {
		if spec == "stdout" {
			handlers = append(handlers, handlerStdoutLogger)
			log.Println("Registered stdout handler.")
		} else if strings.HasPrefix(spec, "traq") {
			traqSpec := strings.Split(spec, ";")
			if len(traqSpec) != 4 {
				return nil, errors.New("traq handler needs exactly 4 arguments")
			}
			if _, err := url.Parse(traqSpec[1]); err != nil {
				return nil, fmt.Errorf("malformed traq origin: %w", err)
			}
			webhookID, err := uuid.FromString(traqSpec[2])
			if err != nil {
				return nil, fmt.Errorf("malformed webhook ID: %w", err)
			}
			handlers = append(handlers, handlerTraqWebhookLogger(traqSpec[1], webhookID, traqSpec[3]))
			log.Printf("Registered traq webhook handler with origin %s.\n", traqSpec[1])
		}
	}
	return
}

func cmdListen() error {
	// parse handlers
	log.Println("Parsing handlers...")
	handlers, err := initHandlers()
	if err != nil {
		return fmt.Errorf("invalid handlers, check HANDLERS env-var syntax: %w", err)
	}
	log.Printf("Successfully initialized total of %d handler(s)!\n", len(handlers))

	// init db
	log.Println("Initializing database...")
	repo, err := gorm.NewGormRepository()
	if err != nil {
		return fmt.Errorf("encountered an error while initializing database: %w", err)
	}
	log.Println("Executing migration...")
	if err := migrate.Migrate(repo.DB()); err != nil {
		return fmt.Errorf("encountered an error while executing database migration: %w", err)
	}
	log.Println("Successfully initialized database connection and executed migrations!")

	// connect
	l := &listener{handlers: handlers, repo: repo}
	log.Println("Connecting to the stream...")
	ctx, cancel := context.WithCancel(context.Background())

	tweetsChan := make(chan *gotwtr.ConnectToStreamResponse, 5)
	stream, err := twitter.NewStream(
		tweetsChan,
		&twitter.ConnectToStreamOption{
			Expansions:  []gotwtr.Expansion{"author_id"},
			TweetFields: []gotwtr.TweetField{"id", "author_id", "entities"},
			UserFields:  []gotwtr.UserField{"name"},
		},
	)

	streamErr := make(chan error)
	go func() {
		if err := stream.StartWithAutoReconnect(ctx); err != nil {
			streamErr <- err
		}
	}()
	go func() {
		for {
			select {
			case event, ok := <-tweetsChan:
				if !ok {
					return
				}
				l.processTweet(event)
			case <-ctx.Done():
				return
			}
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sigChan:
		log.Println("Received SIGTERM or SIGINT signal, closing the stream...")
	case err := <-streamErr:
		cancel()
		return fmt.Errorf("received error from stream, abnormal shutdown: %w", err)
	}

	cancel()
	stream.Stop()
	log.Println("Stream closed. See you next time!")
	return nil
}
