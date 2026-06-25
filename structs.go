package main

import (
	"github.com/kayrabulbul/gator/internal/config"
	"github.com/kayrabulbul/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}
