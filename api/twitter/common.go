package twitter

import (
	"net/http"

	"github.com/sivchari/gotwtr"

	"github.com/motoki317/vrc-world-tweets-fetcher/utils"
)

var c *gotwtr.Client

func init() {
	hc := http.Client{}
	c = gotwtr.New(utils.MustGetEnv("TWITTER_BEARER_TOKEN"), gotwtr.WithHTTPClient(&hc))
}
