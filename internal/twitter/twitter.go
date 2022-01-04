package twitter

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Tweet = twitter.Tweet

type Twitter struct {
	client *twitter.Client
}

func New(ctx context.Context) (*Twitter, error) {
	ckey := os.Getenv("TWITTER_CONSUMER_KEY")
	if ckey == "" {
		return nil, fmt.Errorf("TWITTER_CONSUMER_KEY is not specified")
	}

	csec := os.Getenv("TWITTER_CONSUMER_SECRET")
	if csec == "" {
		return nil, fmt.Errorf("TWITTER_CONSUMER_SECRET is not specified")
	}

	atok := os.Getenv("TWITTER_ACCESS_TOKEN")
	if atok == "" {
		return nil, fmt.Errorf("TWITTER_ACCESS_TOKEN is not specified")
	}

	asec := os.Getenv("TWITTER_ACCESS_SECRET")
	if asec == "" {
		return nil, fmt.Errorf("TWITTER_ACCESS_SECRET is not specified")
	}

	config := oauth1.NewConfig(ckey, csec)
	token := oauth1.NewToken(atok, asec)
	httpClient := config.Client(ctx, token)

	return &Twitter{
		client: twitter.NewClient(httpClient),
	}, nil
}

func (t *Twitter) GetAllTweets() ([]twitter.Tweet, error) {
	tweets := make([]twitter.Tweet, 0, 1000)

	params := &twitter.UserTimelineParams{
		Count:           200,
		TrimUser:        twitter.Bool(false),
		ExcludeReplies:  twitter.Bool(false),
		IncludeRetweets: twitter.Bool(true),
		TweetMode:       "extended",
	}

	var remain int
	for {
		ok, err := func() (bool, error) {
			ts, resp, err := t.client.Timelines.UserTimeline(params)
			if err != nil {
				return false, err
			}
			defer resp.Body.Close()

			if len(ts) == 0 {
				return false, nil
			}
			tweets = append(tweets, ts...)

			maxID := ts[len(ts)-1].ID
			if maxID == 0 {
				return false, fmt.Errorf("max_id is zero")
			}
			params.MaxID = maxID - 1

			remain, err = strconv.Atoi(resp.Header.Get("x-rate-limit-remaining"))
			if err != nil {
				return false, fmt.Errorf("failed to parse x-rate-limit-remaining: %w", err)
			}
			if remain < 1 {
				reset, err := strconv.ParseInt(resp.Header.Get("x-rate-limit-reset"), 10, 64)
				if err != nil {
					return false, fmt.Errorf("failed to parse x-rate-limit-reset: %w", err)
				}
				time.Sleep(time.Until(time.Unix(reset, 0)))
			}

			return true, nil
		}()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}

		// log.Printf("remain: %d, len: %d, max_id: %d", remain, len(tweets), params.MaxID)

		time.Sleep(time.Second)
	}

	return tweets, nil
}

func (t *Twitter) DeleteTweet(tweet twitter.Tweet) error {
	_, _, err := t.client.Statuses.Destroy(tweet.ID, &twitter.StatusDestroyParams{ID: tweet.ID})
	return err
}

func (t *Twitter) DeleteTweetByID(id int64) error {
	_, _, err := t.client.Statuses.Destroy(id, &twitter.StatusDestroyParams{ID: id})
	return err
}
