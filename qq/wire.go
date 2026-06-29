package qq

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// --- feed wire types (hot list + channel feed) ---

type wireHotResp struct {
	Ret    int          `json:"ret"`
	IDList []wireIDList `json:"idlist"`
}

type wireIDList struct {
	IDsHash  string        `json:"ids_hash"`
	HasMore  int           `json:"has_more"`
	NewsList []wireArticle `json:"newslist"`
}

type wireArticle struct {
	ID           string            `json:"id"`
	ArticleType  string            `json:"articletype"`
	Title        string            `json:"title"`
	Abstract     string            `json:"abstract"`
	NLPAbstract  string            `json:"nlpAbstract"`
	SURL         string            `json:"surl"`
	URL          string            `json:"url"`
	Time         string            `json:"time"`
	Timestamp    int64             `json:"timestamp"`
	ChlID        string            `json:"chlid"`
	ChlName      string            `json:"chlname"`
	UserAddress  string            `json:"userAddress"`
	EntityType   string            `json:"entity_type"`
	ReadCount    int               `json:"readCount"`
	LikeInfo     int               `json:"likeInfo"`
	ShareCount   int               `json:"shareCount"`
	CommentNum   int               `json:"commentNum"`
	Comments     int               `json:"comments"`
	CollectCount int               `json:"collect_count"`
	HotEvent     wireHotEvent      `json:"hotEvent"`
	VideoChannel *wireVideoChannel `json:"video_channel"`
}

type wireHotEvent struct {
	ID         string `json:"id"`
	Ranking    int    `json:"ranking"`
	Title      string `json:"title"`
	HotScore   int    `json:"hotScore"`
	IsTop      int    `json:"is_top"`
	SearchWord string `json:"search_word"`
}

type wireVideoChannel struct {
	EGID  string     `json:"egid"`
	Video wireVideo  `json:"video"`
}

type wireVideo struct {
	VID        string `json:"vid"`
	Duration   string `json:"duration"`
	Img        string `json:"img"`
	PlayURL    string `json:"playurl"`
	PlayCount  int    `json:"playcount"`
	ScreenType int    `json:"screenType"`
}

// --- hot words wire types ---

type wireHotWordResp struct {
	Ret      int           `json:"ret"`
	Type     int           `json:"type"`
	TraceID  string        `json:"trace_id"`
	TopWords wireTopWords  `json:"topWords"`
	HotList  []wireArticle `json:"hotlist"`
	TabList  []wireTab     `json:"tablist"`
}

type wireTopWords struct {
	Fixed           []wireWord `json:"fixed"`
	Alternate       []wireWord `json:"alternate"`
	RefreshInterval int        `json:"refreshInterval"`
	DisplayInterval int        `json:"displayInterval"`
}

type wireWord struct {
	Word string `json:"word"`
	From string `json:"from"`
}

type wireTab struct {
	TabName string `json:"tabName"`
	TabDesc string `json:"tabDesc"`
	TabID   string `json:"tabId"`
}

// --- article detail wire types (window.DATA) ---

type wireArticleDetail struct {
	ArticleID       string                     `json:"article_id"`
	ArticleType     string                     `json:"article_type"`
	Title           string                     `json:"title"`
	Desc            string                     `json:"desc"`
	Abstract        string                     `json:"abstract"`
	Media           string                     `json:"media"`
	MediaID         string                     `json:"media_id"`
	PubTime         string                     `json:"pubtime"`
	CommentID       string                     `json:"comment_id"`
	Catalog1        string                     `json:"catalog1"`
	UserAddress     string                     `json:"userAddress"`
	ContentWordsNum int                        `json:"contentWordsNum"`
	ArticleIsPay    bool                       `json:"article_is_pay"`
	Tags            string                     `json:"tags"`
	// Field name in the API is misspelled "interationCount" (missing 'c').
	InterationCount wireInterationCount        `json:"interationCount"`
	OriginContent   wireOriginContent          `json:"originContent"`
	OriginAttribute map[string]json.RawMessage `json:"originAttribute"`
}

type wireInterationCount struct {
	Like    int `json:"like"`
	Collect int `json:"collect"`
	Share   int `json:"share"`
}

type wireOriginContent struct {
	Text    string `json:"text"`
	Version string `json:"version"`
}

type wireImgAsset struct {
	URL        string `json:"url"`
	BigOrigURL string `json:"bigOrigUrl"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Desc       string `json:"desc"`
	IsLong     int    `json:"islong"`
}

type wireVideoAsset struct {
	VID      string `json:"vid"`
	Img      string `json:"img"`
	Duration string `json:"duration"`
	PlayURL  string `json:"playurl"`
	Title    string `json:"title"`
}

// --- converters ---

// articleFrom converts a feed wire item to a public Article.
// rank is the 1-based position in the output list (after skipping the type-560 header).
func articleFrom(w wireArticle, rank int) Article {
	abs := w.Abstract
	if abs == "" {
		abs = w.NLPAbstract
	}
	u := w.SURL
	if u == "" {
		u = w.URL
	}
	comments := w.CommentNum
	if comments == 0 {
		comments = w.Comments
	}
	a := Article{
		Rank:        rank,
		ID:          w.ID,
		Title:       w.Title,
		Abstract:    abs,
		URL:         u,
		Time:        w.Time,
		Timestamp:   w.Timestamp,
		Publisher:   w.ChlName,
		PublisherID: w.ChlID,
		UserAddress: w.UserAddress,
		EntityType:  w.EntityType,
		HotScore:    w.HotEvent.HotScore,
		ReadCount:   w.ReadCount,
		Likes:       w.LikeInfo,
		Shares:      w.ShareCount,
		Comments:    comments,
	}
	if w.VideoChannel != nil {
		v := w.VideoChannel.Video
		a.VideoDuration = v.Duration
		a.VideoVID = v.VID
		a.VideoPlayURL = v.PlayURL
		a.VideoPlays = v.PlayCount
	}
	return a
}

var (
	reComment   = regexp.MustCompile(`<!--.*?-->`)
	reVertCard  = regexp.MustCompile(`VERTICAL_CARD_(?:BEGIN|END)_\d+`)
	reHTMLTag   = regexp.MustCompile(`<[^>]*>`)
	reMultiLine = regexp.MustCompile(`\n{3,}`)
)

// stripHTML removes HTML tags, placeholder comments, and collapses whitespace.
func stripHTML(s string) string {
	s = reComment.ReplaceAllString(s, "")
	s = reVertCard.ReplaceAllString(s, "")
	s = reHTMLTag.ReplaceAllString(s, "")
	// Decode common HTML entities
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&quot;", `"`)
	s = reMultiLine.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}

// articleDetailFrom converts a wireArticleDetail to a public ArticleDetail.
// articleID is the canonical ID from the feed (the detail page may truncate it for videos).
// withBody controls whether originContent.text is stripped and included.
func articleDetailFrom(w wireArticleDetail, articleID string, withBody bool) ArticleDetail {
	abs := w.Desc
	if abs == "" {
		abs = w.Abstract
	}
	d := ArticleDetail{
		ID:          articleID,
		Title:       w.Title,
		Abstract:    abs,
		URL:         "https://view.inews.qq.com/a/" + articleID,
		Publisher:   w.Media,
		PublisherID: w.MediaID,
		PubTime:     w.PubTime,
		Category:    w.Catalog1,
		WordCount:   w.ContentWordsNum,
		UserAddress: w.UserAddress,
		Likes:       w.InterationCount.Like,
		Collects:    w.InterationCount.Collect,
		Shares:      w.InterationCount.Share,
		IsPaid:      w.ArticleIsPay,
		Tags:        w.Tags,
	}
	if withBody && w.OriginContent.Text != "" {
		d.Body = stripHTML(w.OriginContent.Text)
	}
	d.Images = extractImages(w.OriginAttribute)
	d.Videos = extractVideos(w.OriginAttribute)
	return d
}

// numberedKey returns the numeric suffix of keys like "IMG_0", "VIDEO_2".
func numberedKey(prefix, key string) (int, bool) {
	s := strings.TrimPrefix(key, prefix)
	if s == key {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}

func extractImages(attr map[string]json.RawMessage) []ArticleImage {
	type indexed struct {
		n int
		v ArticleImage
	}
	var items []indexed
	for k, raw := range attr {
		n, ok := numberedKey("IMG_", k)
		if !ok {
			continue
		}
		var w wireImgAsset
		if err := json.Unmarshal(raw, &w); err != nil {
			continue
		}
		u := w.BigOrigURL
		if u == "" {
			u = w.URL
		}
		items = append(items, indexed{n, ArticleImage{
			URL:     u,
			Width:   w.Width,
			Height:  w.Height,
			Caption: w.Desc,
			IsLong:  w.IsLong == 1,
		}})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].n < items[j].n })
	out := make([]ArticleImage, len(items))
	for i, it := range items {
		out[i] = it.v
	}
	return out
}

func extractVideos(attr map[string]json.RawMessage) []ArticleVideo {
	type indexed struct {
		n int
		v ArticleVideo
	}
	var items []indexed
	for k, raw := range attr {
		n, ok := numberedKey("VIDEO_", k)
		if !ok {
			continue
		}
		var w wireVideoAsset
		if err := json.Unmarshal(raw, &w); err != nil {
			continue
		}
		items = append(items, indexed{n, ArticleVideo{
			VID:      w.VID,
			Cover:    w.Img,
			Duration: w.Duration,
			PlayURL:  w.PlayURL,
			Title:    w.Title,
		}})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].n < items[j].n })
	out := make([]ArticleVideo, len(items))
	for i, it := range items {
		out[i] = it.v
	}
	return out
}

// hotWordFrom converts an alternate word entry to a HotWord.
func hotWordFrom(w wireWord, rank int) HotWord {
	return HotWord{Rank: rank, Word: w.Word}
}

// parseWindowData extracts the wireArticleDetail from an HTML page's window.DATA script.
var reWindowData = regexp.MustCompile(`(?s)window\.DATA\s*=\s*(\{.+?\})\s*;`)

func parseWindowData(html string) (*wireArticleDetail, error) {
	m := reWindowData.FindStringSubmatch(html)
	if len(m) < 2 {
		return nil, fmt.Errorf("window.DATA not found in page")
	}
	var d wireArticleDetail
	if err := json.Unmarshal([]byte(m[1]), &d); err != nil {
		return nil, fmt.Errorf("parse window.DATA: %w", err)
	}
	return &d, nil
}
