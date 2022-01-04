package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shuymn/ephmr/internal/twitter"
)

const TWEET_FILE_PATH = "./tweet.json"

func main() {
	b, err := os.ReadFile(TWEET_FILE_PATH)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open tweet.json: %w", err))
	}

	var tweets []Tweet
	if err = json.Unmarshal(b, &tweets); err != nil {
		log.Fatal(fmt.Errorf("failed to unmarshal tweet.json: %w", err))
	}

	ctx := context.Background()
	tw, err := twitter.New(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize twitter client: %w", err))
	}

	total := len(tweets)
	for i, tweet := range tweets {
		id, err := strconv.ParseInt(tweet.Tweet.ID, 10, 64)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to parse tweet id: %w", err))
		}

		log.Printf("[%d] delete(%d%%) => id: %s (%s)", i, 100*i/total, tweet.Tweet.ID, tweet.Tweet.CreatedAt)

		if err = tw.DeleteTweetByID(id); err != nil {
			if err.Error() == "twitter: 144 No status found with that ID." {
				continue
			}
			log.Fatal(fmt.Errorf("failed to delete tweet: %w", err))
		}
	}
}
