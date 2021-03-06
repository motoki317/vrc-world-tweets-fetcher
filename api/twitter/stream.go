package twitter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/samber/lo"
	"github.com/sivchari/gotwtr"
)

// Down from here includes modified code from https://github.com/sivchari/gotwtr/blob/1a979e230d898bf6c8bd5c64a8ae47b27f2d3a4f/filtered_stream.go

const (
	connectToStreamURL = "https://api.twitter.com/2/tweets/search/stream"
)

func stopped(done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

type ConnectToStreamOption struct {
	Expansions  []gotwtr.Expansion
	MediaFields []gotwtr.MediaField
	PlaceFields []gotwtr.PlaceField
	PollFields  []gotwtr.PollField
	TweetFields []gotwtr.TweetField
	UserFields  []gotwtr.UserField
}

func join[T ~string](elems []T, sep string) string {
	return strings.Join(lo.Map[T, string](elems, func(elem T, _ int) string { return string(elem) }), sep)
}

func (t *ConnectToStreamOption) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(t.Expansions) > 0 {
		q.Add("expansions", join(t.Expansions, ","))
	}
	if len(t.MediaFields) > 0 {
		q.Add("media.fields", join(t.MediaFields, ","))
	}
	if len(t.PlaceFields) > 0 {
		q.Add("place.fields", join(t.PlaceFields, ","))
	}
	if len(t.PollFields) > 0 {
		q.Add("poll.fields", join(t.PollFields, ","))
	}
	if len(t.TweetFields) > 0 {
		q.Add("tweet.fields", join(t.TweetFields, ","))
	}
	if len(t.UserFields) > 0 {
		q.Add("user.fields", join(t.UserFields, ","))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

type Stream struct {
	client *http.Client
	option *ConnectToStreamOption

	tweets chan<- *gotwtr.ConnectToStreamResponse

	done chan struct{}
	wg   *sync.WaitGroup
}

// StartWithAutoReconnect connects to stream, and reconnects to the stream as long as no errors are caught.
// Blocks on success.
func (s *Stream) StartWithAutoReconnect(ctx context.Context) error {
	if stopped(s.done) {
		return nil
	}
	for {
		if err := s.Start(ctx); err != nil {
			return err
		}
		if stopped(s.done) {
			return nil
		}
		log.Println("Reconnecting to tweet stream...")
	}
}

func (s *Stream) readLoop(dec *json.Decoder) error {
	type decodeResult struct {
		resp *gotwtr.ConnectToStreamResponse
		err  error
	}
	decodeChan := make(chan decodeResult)
	for {
		go func() {
			var response gotwtr.ConnectToStreamResponse
			err := dec.Decode(&response)
			decodeChan <- decodeResult{resp: &response, err: err}
		}()
		select {
		case res := <-decodeChan:
			if res.err != nil {
				if res.err == io.EOF {
					return nil
				}
				return res.err
			}
			s.tweets <- res.resp
		case <-s.done:
			return nil
		}
	}
}

// Start connects to stream, and blocks on success.
func (s *Stream) Start(ctx context.Context) error {
	s.wg.Add(1)
	defer s.wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, connectToStreamURL, nil)
	if err != nil {
		return fmt.Errorf("connect to stream new request with ctx: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	s.option.addQuery(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &gotwtr.HTTPError{
			APIName: "connect to stream",
			Status:  resp.Status,
			URL:     req.URL.String(),
		}
	}

	log.Println("Connected! Now receiving new tweets...")
	return s.readLoop(json.NewDecoder(resp.Body))
}

func (s *Stream) Stop() {
	close(s.done)
	s.wg.Wait()
}

func NewStream(tweets chan<- *gotwtr.ConnectToStreamResponse, opt ...*ConnectToStreamOption) (*Stream, error) {
	var option ConnectToStreamOption
	switch len(opt) {
	case 0:
		// do nothing
	case 1:
		option = *opt[0]
	default:
		return nil, errors.New("connect to stream: only one option is allowed")
	}

	return &Stream{
		client: hc,
		option: &option,
		tweets: tweets,
		done:   make(chan struct{}),
		wg:     &sync.WaitGroup{},
	}, nil
}
