package twitter

import (
	"net/http"

	"github.com/sivchari/gotwtr"

	"github.com/motoki317/vrc-world-tweets-fetcher/utils"
)

var (
	hc    *http.Client
	c     *gotwtr.Client
	token = utils.MustGetEnv("TWITTER_BEARER_TOKEN")
)

func init() {
	hc = &http.Client{}
	c = gotwtr.New(token, gotwtr.WithHTTPClient(hc))
}
