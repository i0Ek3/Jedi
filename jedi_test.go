package main

import (
    "encoding/json"
    "os"
    "testing"

    "github.com/i0Ek3/asrt"
)

var (
    d   defaultMatcher
    r   rssMatcher
    err error

    feed       *Feed
    results    []*Result
)

const (
    searchTerm = "economic"
)

// TODO: 
// finish basic test process but without data filling, 
// next step should update date from the code.

func TestDefaultSearch(t *testing.T) {
    got, _ := d.Search(feed, searchTerm)
    want := feed
    asrt.Equal(got, want)
}

func TestRSSSearch(t *testing.T) {
    feed = &Feed{
        Name: "foxnews",
        URI:  "https://feeds.foxnews.com/foxnews/national?format=rss",
        Type: "rss",
    }
    results, _ := r.Search(feed, searchTerm)
    ShowMe(results)
    want := []*Result{}
    asrt.Equal(results, want)
}

func TestRetrieve(t *testing.T) {
    file, _ := os.Open("./data.json")
    defer file.Close()
    _ = json.NewDecoder(file).Decode(feed)

    got, _ := RetrieveFeeds()
    want := feed
    asrt.Equal(got, want)
}
