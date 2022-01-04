package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shuymn/ephmr/internal/twitter"
)

func main() {
	ctx := context.Background()
	twtr, err := twitter.New(ctx)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize twitter client: %w", err))
	}

	tweets, err := twtr.GetAllTweets()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get all tweets: %w", err))
	}

	total := len(tweets)
	for i, tweet := range tweets {
		if isArchive(tweet) {
			log.Printf("delete(%d%%) => id: %d (%s)", 100*i/total, tweet.ID, tweet.CreatedAt)

			if err = twtr.DeleteTweet(tweet); err != nil {
				log.Fatal(fmt.Errorf("failed to delete tweet: %w", err))
			}

			time.Sleep(time.Second)
		}
	}
}

func isArchive(tweet twitter.Tweet) bool {
	createdAt, err := tweet.CreatedAtTime()
	if err != nil {
		log.Print(err)
		return false
	}
	return createdAt.Before(time.Now().Add(-1 * 3 * 24 * time.Hour))
}
