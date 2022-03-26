package utils

import (
	"os"
)

func MustGetEnv(name string) string {
	if s := os.Getenv(name); s != "" {
		return s
	}
	panic("Environment variable " + name + " is not defined")
}
