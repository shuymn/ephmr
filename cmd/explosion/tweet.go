package main

type Tweet struct {
	Tweet TweetDetail `json:"tweet"`
}

type TweetDetail struct {
	Retweeted        bool     `json:"retweeted"`
	Source           string   `json:"source"`
	Entities         Entities `json:"entities"`
	DisplayTextRange []string `json:"display_text_range"`
	FavoriteCount    string   `json:"favorite_count"`
	IDStr            string   `json:"id_str"`
	Truncated        bool     `json:"truncated"`
	RetweetCount     string   `json:"retweet_count"`
	ID               string   `json:"id"`
	CreatedAt        string   `json:"created_at"`
	Favorited        bool     `json:"favorited"`
	FullText         string   `json:"full_text"`
	Lang             string   `json:"lang"`
}

type Entities struct {
	Hashtags     []interface{} `json:"hashtags"`
	Symbols      []interface{} `json:"symbols"`
	UserMentions []interface{} `json:"user_mentions"`
	Urls         []interface{} `json:"urls"`
}
