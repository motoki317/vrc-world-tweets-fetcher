package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds | log.Lmsgprefix)
	log.SetPrefix("[vrc-world-tweets-fetcher] ")

	flag.Parse()
}

const usage = `Usage: %s [init|listen]

init:	Initializes stream rules. Deletes any existing rules if found.

listen:	Connects to stream and start receiving tweets.
	Use HANDLERS environment variable to specify handlers.

	HANDLERS:
		Comma-separated list of handlers on newly found world.
		Defaults to "stdout". Multiple instances of the same kind of handlers allowed.

		Allowed values:
			- stdout: Logs world info to stdout.
			- traq: Logs world info to traQ webhook. Syntax: traq;ORIGIN;WEBHOOK_ID;WEBHOOK_SECRET, example: traq;https://q.trap.jp;00000000-0000-0000-0000-000000000000;my-secret

		Example: stdout,traq:https://q.trap.jp;00000000-0000-0000-0000-000000000000;my-secret,traq;https://q.toki317.dev;00000000-0000-0000-0000-000000000000;my-secret-2`

func main() {
	args := flag.Args()
	if len(args) == 0 {
		log.Println(fmt.Sprintf(usage, os.Args[0]))
		return
	}

	var cmdHandler func() error
	switch args[0] {
	case "init":
		cmdHandler = cmdInit
	case "listen":
		cmdHandler = cmdListen
	default:
		log.Println(fmt.Sprintf(`Usage: %s [init|listen]`, os.Args[0]))
		os.Exit(1)
	}

	if err := cmdHandler(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
