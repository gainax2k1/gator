package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/gainax2k1/gator/internal/database"
	"github.com/lib/pq"
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

		var postParams database.CreatePostParams
		postParams.Title = rssitems.Title
		postParams.Url = rssitems.Link
		postParams.FeedID = feed.ID

		DBpost, err := s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			// need to check if err is url already exists, and ignore it

			pqErrorVal, ok := err.(*pq.Error)
			if ok {
				if pqErrorVal.Code == "23505" {
					fmt.Println("23505 hit")
					// 23505 is already exists error, prolly URL, expected and ignored
					continue
				}
			}
			fmt.Printf("error creating post in DB: %v\n", err)
			continue
		}
		fmt.Printf("created post: %v", DBpost.Title)

		if rssitems.PubDate != "" {
			var updatePublishedDateParams database.UpdatePublishedAtParams
			updatePublishedDateParams.ID = DBpost.ID

			//need to evaluate pubdate in terms of how to make it
			//work with Time

			fixedTime, err := StandardizeTime(rssitems.PubDate)
			if err != nil {
				//if time was unable to parse correctly, fixed time should hold
				// zero value ( time.Time{} ), so....?
				fmt.Println("\n%w\n", err)
			} else {

				updatePublishedDateParams.PublishedAt.Time = fixedTime // correct
				updatePublishedDateParams.PublishedAt.Valid = true
				s.db.UpdatePublishedAt(context.Background(), updatePublishedDateParams)
			}
		}
		if rssitems.Description != "" {
			var updateDescriptionParams database.UpdateDescriptionParams
			updateDescriptionParams.Description = rssitems.Description
			updateDescriptionParams.ID = DBpost.ID
			_, err = s.db.UpdateDescription(context.Background(), updateDescriptionParams)
			if err != nil {
				return fmt.Errorf("error updating description: %v", err)
			}
		}

		fmt.Printf("dbpost: %v", DBpost) // for troubleshooting
	}

	return nil
}

func PrintPost(post database.Post) error {
	fmt.Printf("Title: %s\n", post.Title)
	fmt.Printf("Description: %s\n", post.Description)
	fmt.Printf("Published: %v\n", post.PublishedAt.Time)
	fmt.Printf("URL: %s\n", post.Url)
	return nil
}
