# vrc-world-tweets-fetcher

This simple app streams tweets from Twitter, parses its content, and if it contained VRChat world link,
posts the link to several media.

Post destinations are fully configurable.

## Usage

```shell
$ ./vrc-world-tweets-fetcher 
Usage: ./vrc-world-tweets-fetcher [init|listen]

init:   Initializes stream rules. Deletes any existing rules if found.

listen: Connects to stream and start receiving tweets.
        Use HANDLERS environment variable to specify handlers.

        HANDLERS:
                Comma-separated list of handlers on newly found world.
                Defaults to "stdout". Multiple instances of the same kind of handlers allowed.

                Allowed values:
                        - stdout: Logs world info to stdout.
                        - traq: Logs world info to traQ webhook. Syntax: traq;ORIGIN;WEBHOOK_ID;WEBHOOK_SECRET, example: traq;https://q.trap.jp;00000000-0000-0000-0000-000000000000;my-secret

                Example: stdout,traq:https://q.trap.jp;00000000-0000-0000-0000-000000000000;my-secret,traq;https://q.toki317.dev;00000000-0000-0000-0000-000000000000;my-secret-2
```

## Links

- [Twitter API Documentation | Docs | Twitter Developer Platform](https://developer.twitter.com/en/docs/twitter-api)
- [VRChat API Documentation](https://vrchatapi.github.io/)
- [traPtitech/traQ: traQ - traP Internal Messenger Application Backend](https://github.com/traPtitech/traQ)
