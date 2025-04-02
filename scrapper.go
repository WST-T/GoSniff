package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/WST-T/GoSniff/internal/database"
	"github.com/google/uuid"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Starting scraper with concurrency: %d and time between requests: %s", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error getting feeds to fetch: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	// Validate URL before attempting to fetch
	if feed.Url == "" {
		log.Printf("Error: Empty URL for feed ID %s, name: %s", feed.ID, feed.Name)
		markFeedWithError(db, feed, "Empty URL")
		return
	}

	// Validate URL format
	_, err := url.ParseRequestURI(feed.Url)
	if err != nil {
		log.Printf("Error: Invalid URL format for feed %s (%s): %v", feed.Name, feed.Url, err)
		markFeedWithError(db, feed, fmt.Sprintf("Invalid URL: %v", err))
		return
	}

	_, err = db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed %s (%s): %v", feed.Name, feed.Url, err)
		// Continue with the rest of the process despite the error
		return
	}

	for _, item := range rssFeed.Channel.Item {
		publishedAt, err := item.ParsedPublishedAt()
		if err != nil {
			log.Printf("Couldn't parse date %v: %v", item.PubDate, err)
			continue
		}

		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Failed to create post: %v", err)
			continue
		}
		log.Printf("Found post %v on feed %v", item.Title, feed.Name)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

// Helper function to mark a feed with error information
// Note: You'll need to add a column to the feeds table and update the database queries
// to implement this fully, but this function shows the concept
func markFeedWithError(db *database.Queries, feed database.Feed, errorMsg string) {
	// Log the error message for debugging purposes
	log.Printf("Feed error for %s (%s): %s", feed.Name, feed.ID, errorMsg)

	// This is a placeholder for future implementation
	// Ideally, we would update the feed with error information
	// For now, we just mark it as fetched to avoid constant retries
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking problematic feed as fetched: %v", err)
	}
}
