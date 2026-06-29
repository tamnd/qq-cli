package qq

import "errors"

// ErrNotFound is returned when the API returns no result for a known-valid ID.
var ErrNotFound = errors.New("not found")

// Article is a news item from the hot list or a channel feed.
// Used by: hot, channel.
type Article struct {
	Rank        int    `json:"rank"         table:",right"`
	ID          string `json:"id"`
	Title       string `json:"title"        table:",truncate"`
	Abstract    string `json:"abstract"     table:"-"`
	URL         string `json:"url"          kit:"url" table:",truncate"`
	Time        string `json:"time"`
	Timestamp   int64  `json:"timestamp"    table:"-"`
	Publisher   string `json:"publisher"`
	PublisherID string `json:"publisher_id" table:"-"`
	UserAddress string `json:"user_address" table:"-"`
	EntityType  string `json:"entity_type"`
	HotScore    int    `json:"hot_score"    table:",right"`
	ReadCount   int    `json:"read_count"   table:"-"`
	Likes       int    `json:"likes"        table:"-"`
	Shares      int    `json:"shares"       table:"-"`
	Comments    int    `json:"comments"     table:"-"`
	// Video-only (zero/empty for text articles)
	VideoDuration string `json:"video_duration,omitempty" table:"-"`
	VideoVID      string `json:"video_vid,omitempty"      table:"-"`
	VideoPlayURL  string `json:"video_play_url,omitempty" table:"-"`
	VideoPlays    int    `json:"video_plays,omitempty"    table:"-"`
}

// ArticleDetail is a fully fetched article with body content.
// Used by: article.
type ArticleDetail struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"        table:",truncate"`
	Abstract    string         `json:"abstract"     table:",truncate"`
	URL         string         `json:"url"          kit:"url" table:",truncate"`
	Publisher   string         `json:"publisher"`
	PublisherID string         `json:"publisher_id" table:"-"`
	PubTime     string         `json:"pub_time"`
	Category    string         `json:"category"     table:"-"`
	WordCount   int            `json:"word_count"   table:"-"`
	UserAddress string         `json:"user_address" table:"-"`
	Likes       int            `json:"likes"        table:"-"`
	Collects    int            `json:"collects"     table:"-"`
	Shares      int            `json:"shares"       table:"-"`
	IsPaid      bool           `json:"is_paid"      table:"-"`
	Tags        string         `json:"tags"         table:"-"`
	Body        string         `json:"body,omitempty"   table:"-"`
	Images      []ArticleImage `json:"images,omitempty" table:"-"`
	Videos      []ArticleVideo `json:"videos,omitempty" table:"-"`
}

// ArticleImage is one image embedded in an article body.
type ArticleImage struct {
	URL     string `json:"url"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Caption string `json:"caption"`
	IsLong  bool   `json:"is_long"`
}

// ArticleVideo is one video embedded in an article body.
type ArticleVideo struct {
	VID      string `json:"vid"`
	Cover    string `json:"cover"`
	Duration string `json:"duration"`
	PlayURL  string `json:"play_url"`
	Title    string `json:"title"`
}

// HotWord is one trending search word suggestion.
// Used by: trending.
type HotWord struct {
	Rank int    `json:"rank" table:",right"`
	Word string `json:"word" kit:"url" table:",truncate"`
}

// Channel is one entry in the QQ News channel directory.
// Used by: channels.
type Channel struct {
	ChannelID string `json:"channel_id"`
	NameCN    string `json:"name_cn"`
	NameEN    string `json:"name_en"`
	IsLocal   bool   `json:"is_local"  table:"-"`
	HasFeed   bool   `json:"has_feed"`
	Link      string `json:"link,omitempty" table:"-"`
}
