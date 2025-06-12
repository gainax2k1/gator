package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var rssfeed *RSSFeed = &RSSFeed{}
	var readerbody io.Reader

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, readerbody)
	if err != nil {
		return rssfeed, err
	}

	req.Header.Set("User-Agent", "gator") // identify self to server

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return rssfeed, fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return rssfeed, fmt.Errorf("error reading response body: %w", err)
	}

	if err = xml.Unmarshal(data, rssfeed); err != nil {
		return rssfeed, err
	}

	// CLEAN rssfeed HERE
	rssfeed.Channel.Title = html.UnescapeString(rssfeed.Channel.Title)
	rssfeed.Channel.Description = html.UnescapeString((rssfeed.Channel.Description))
	for i, rssitem := range rssfeed.Channel.Item {
		rssfeed.Channel.Item[i].Title = html.UnescapeString((rssitem.Title))
		rssfeed.Channel.Item[i].Description = html.UnescapeString((rssitem.Description))
	}

	return rssfeed, nil

}

func scrapeFeeds(s *state) error {

	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("error marking feed [%s] fetched: %w", feed.Name, err)
	}

	RSSFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed [%s]: %w", feed.Name, err)
	}

	fmt.Printf("RSS Channel: %s\n", RSSFeed.Channel.Title)
	fmt.Printf("Description: %s\n", RSSFeed.Channel.Description)

	for _, rssitems := range RSSFeed.Channel.Item {
		fmt.Printf("- Title: %s\n", rssitems.Title)
		//fmt.Println(" - Link: %s", rssitems.Link)
		//fmt.Println(" - Description: %s", rssitems.Description)
		//fmt.Println(" - Published: %s", rssitems.PubDate)

	}

	return nil
}
