package main

import (
    "encoding/json"
    "encoding/xml"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "flag"
    "regexp"
    "sync"
    "time"
    _ "net/http/pprof"

    pbar "github.com/schollz/progressbar/v3"
    log "github.com/sirupsen/logrus"
    "github.com/i0Ek3/color"
    "github.com/i0Ek3/noerr"
)

const (
    dataFile   = "data.json"
    searchConf = "search.conf"
    strTerm    = "China"
    numLimit   = 20
)

type (
    // T ready for generic support
    T interface{}

    // item defines the fields associated with the item tag in the rss document
    item struct {
        XMLName     xml.Name `xml:"item"`
        PubDate     string   `xml:"pubDate"`
        Title       string   `xml:"title"`
        Description string   `xml:"description"`
        Link        string   `xml:"link"`
        GUID        string   `xml:"guid"`
        GeoRssPoint string   `xml:"georss:point"`
    }

    // image defines the fields associated with the image tag in the rss document
    image struct {
        XMLName xml.Name `xml:"image"`
        URL     string   `xml:"url"`
        Title   string   `xml:"title"`
        Link    string   `xml:"link"`
    }

    // channel defines the fields associated with the channel tag in the rss document
    channel struct {
        XMLName        xml.Name `xml:"channel"`
        Title          string   `xml:"title"`
        Description    string   `xml:"description"`
        Link           string   `xml:"link"`
        PubDate        string   `xml:"pubDate"`
        LastBuildDate  string   `xml:"lastBuildDate"`
        TTL            string   `xml:"ttl"`
        Language       string   `xml:"language"`
        ManagingEditor string   `xml:"managingEditor"`
        WebMaster      string   `xml:"webMaster"`
        Image          image    `xml:"image"`
        Item           []item   `xml:"item"`
    }

    // rssDocument defines the fields associated with the rss document
    rssDocument struct {
        XMLName xml.Name `xml:"rss"`
        Channel channel  `xml:"channel"`
    }
)

// Result defines the result of a search
type Result struct {
    Field   string
    Content string
}

// Feed defines the struct of a feed 
type Feed struct {
    Name string `json:"site"`
    URI  string `json:"link"`
    Type string `json:"type"`
}

// Matcher defines the Search method
type Matcher interface {
    Search(feed *Feed, searchTerm ...string) ([]*Result, error)
}

type (
    // rssMatcher implements the Matcher interface
    rssMatcher struct{}

    // defaultMatcher implements the default matcher
    defaultMatcher struct{}
)

var (
    matchers = make(map[string]Matcher)
)

// Search returns nil for default matcher
func (d defaultMatcher) Search(feed *Feed, searchTerm ...string) ([]*Result, error) {
    return nil, nil
}

func ShowMe(v T) {
    sign := "+++++++++++++++++++++++++++++++"
    color.Colorize("white", sign)
    log.Println(v)
    color.Colorize("white", sign)
}

// Search finds the specified search term
func (r rssMatcher) Search(feed *Feed, searchTerm ...string) ([]*Result, error) {
    var results []*Result

    str := fmt.Sprintf("Search Feed Type[%s] Site[%s] for URI[%s]\n", feed.Type, feed.Name, feed.URI)
    color.Colorize("blue", str)

    document, err := r.retrieve(feed)
    if err != nil {
        return nil, err
    }

    for _, channelItem := range document.Channel.Item {
        matched, err := regexp.MatchString(searchTerm[0], channelItem.Title)
        noerr.NoErr(err)

        // if matched then save to the result
        if matched {
            results = append(results, &Result{
                Field:   color.Blue("Title"),
                Content: channelItem.Title,
            })
        }

        matched, err = regexp.MatchString(searchTerm[0], channelItem.Description)
        noerr.NoErr(err)

        if matched {
            results = append(results, &Result{
                Field:   color.Blue("Description"),
                Content: channelItem.Description,
            })
        }
    }
    return results, nil
}

// retrieve sends a http get request to fetch feeds then decodes the results
func (r rssMatcher) retrieve(feed *Feed) (*rssDocument, error) {
    if feed.URI == "" {
        return nil, errors.New(color.Red("RSS feed URI required!"))
    }

    resp, err := http.Get(feed.URI)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf(color.Red("HTTP Response Error %d\n"), resp.StatusCode)
    }

    var document rssDocument
    err = xml.NewDecoder(resp.Body).Decode(&document)
    return &document, err
}

// RetrieveFeeds operates the feed data
func RetrieveFeeds() ([]*Feed, error) {
    file, err := os.Open(dataFile)
    noerr.NoErr(err)
    defer file.Close()

    var feeds []*Feed
    err = json.NewDecoder(file).Decode(&feeds)

    return feeds, err
}

// Match searches the result concurrently
func Match(m Matcher, feed *Feed, results chan<- *Result, searchTerm ...string) {
    searchResults, err := m.Search(feed, searchTerm...)
    if err != nil {
        log.Println(err)
        return
    }

    for _, result := range searchResults {
        results <- result
    }
}

// Display displays the result
func Display(limit int, results chan *Result) {
    cnt := 0
    for result := range results {
        cnt += 1
        if cnt < limit {
            fmt.Printf("%s\n%s\n\n", color.Blue(result.Field), color.Cyan(result.Content))
        } else {
            break
        }
    }
}

// Run defines the logic of search
func Process(limit int, searchTerm ...string) {
    feeds, err := RetrieveFeeds()
    noerr.NoErrln(err)

    results := make(chan *Result)

    var wg sync.WaitGroup
    wg.Add(len(feeds))

    for _, feed := range feeds {
        matcher, exists := matchers[feed.Type]
        if !exists {
            matcher = matchers["default"]
        }

        go func(matcher Matcher, feed *Feed) {
            Match(matcher, feed, results, searchTerm...)
            wg.Done()
        }(matcher, feed)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    Display(limit, results)
}

// Register registers a matcher
func Register(feedType string, matcher Matcher) {
    if _, exists := matchers[feedType]; exists {
        fmt.Println(color.Red(feedType, "Matcher already registered!"))
    }

    color.Colorize("blue", "Register", feedType, "matcher")
    matchers[feedType] = matcher
}

// loadFile loads file and excute search
func loadFile(filename string) {
    f, err := os.Open(filename)
    noerr.NoErrln(err)
    defer f.Close()

    searchItem, err := ioutil.ReadAll(f)
    noerr.NoErrln(err)
    setBar(100)
    Process(numLimit, string(searchItem))
}

// init initializes the log output and register service
func init() {
    var m1 rssMatcher
    tag1 := "rss"
    Register(tag1, m1)

    var m2 defaultMatcher
    tag2 := "default"
    Register(tag2, m2)
}

func logSetting() {
    log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
    log.SetOutput(os.Stdout)
    log.SetLevel(log.InfoLevel)
}

// TODO: setBar can be enhenced
func setBar(num int64) {
    bar := pbar.Default(num)
    for i := 0; i < int(num); i++ {
        bar.Add(1)
        time.Sleep(10 * time.Millisecond)
    }
}

var (
    keyword *string
    limit   *int
)

func init() {
    keyword = flag.String("keyword", strTerm, "specific a keyword to query")
    limit = flag.Int("limit", numLimit, "how many matched items shows here")
}

func main() {
    logSetting()

    flag.Parse()

    Process(*limit, *keyword)
    //loadFile(searchConf)

    _ = http.ListenAndServe("0.0.0.0:8899", nil)
}
