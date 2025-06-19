package main

import (
	"fmt"
	"time"
)

func StandardizeTime(pubDateString string) (time.Time, error) {

	/*

		RFC1123: "Mon, 02 Jan 2006 15:04:05 MST"
		RFC1123Z: "Mon, 02 Jan 2006 15:04:05 -0700"
		RFC3339 (ISO 8601): "2006-01-02T15:04:05Z07:00"
	*/
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}
	var t time.Time
	var err error

	for _, layout := range layouts {
		t, err = time.Parse(layout, pubDateString)
		if err == nil {
			break
		}
	}
	if err != nil {
		// Unable to parse the date -- handle gracefully
		return time.Time{}, fmt.Errorf("unable to parse date: %w", err)
	}

	return t, nil

}
