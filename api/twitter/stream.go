package twitter

import (
	"context"

	"github.com/sivchari/gotwtr"
)

func ConnectToStream(ctx context.Context, tweetsChan chan<- gotwtr.ConnectToStreamResponse, errChan chan<- error, opt ...*gotwtr.ConnectToStreamOption) *gotwtr.ConnectToStream {
	return c.ConnectToStream(ctx, tweetsChan, errChan, opt...)
}
