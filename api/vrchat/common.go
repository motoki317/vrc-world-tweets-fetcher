package vrchat

import (
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

const (
	vrchatAPIBasePath = "https://vrchat.com/api/1"
	vrchatAPIKey      = "JlE5Jldo5Jibnk5O5hTx6XVqsJu4WJ26"
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

const userAgent = "vrc-world-tweets-fetcher 1.0; github.com/motoki317/vrc-world-tweets-fetcher"
