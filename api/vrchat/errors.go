package vrchat

import (
	"errors"
)

var (
	ErrServer = errors.New("server returned 5xx response")
	ErrClient = errors.New("server returned 4xx response")
)
