package qq_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/qq-cli/qq"
)

func newTestClient(ts *httptest.Server) *qq.Client {
	cfg := qq.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.ArticleBaseURL = ts.URL
	cfg.Rate = 0
	cfg.Retries = 0
	return qq.NewClient(cfg)
}

// --- fixtures ---

const fixtureHot = `{
  "ret": 0,
  "idlist": [{
    "ids_hash": "abc",
    "has_more": 1,
    "newslist": [
      {
        "id": "TIP2022042216544300",
        "articletype": "560",
        "title": "Header Placeholder",
        "surl": "https://news.qq.com/omn/20220422/a560.htm"
      },
      {
        "id": "20260615A02AHC00",
        "articletype": "0",
        "title": "AI在新闻领域的最新突破",
        "abstract": "人工智能正改变新闻生产方式",
        "nlpAbstract": "",
        "surl": "https://view.inews.qq.com/a/20260615A02AHC00",
        "time": "2026-06-15 10:00:00",
        "timestamp": 1749945600,
        "chlid": "news_news_tech",
        "chlname": "科技",
        "userAddress": "北京",
        "entity_type": "article",
        "readCount": 42000,
        "likeInfo": 1200,
        "shareCount": 300,
        "commentNum": 88,
        "hotEvent": {
          "id": "e1",
          "ranking": 1,
          "title": "AI新闻",
          "hotScore": 9500,
          "is_top": 0,
          "search_word": "AI"
        }
      },
      {
        "id": "20260615B00XYZ00",
        "articletype": "0",
        "title": "体育新闻",
        "abstract": "今日体育头条",
        "surl": "https://view.inews.qq.com/a/20260615B00XYZ00",
        "time": "2026-06-15 11:00:00",
        "timestamp": 1749949200,
        "chlid": "news_news_sports",
        "chlname": "体育",
        "readCount": 21000,
        "likeInfo": 400,
        "shareCount": 100,
        "commentNum": 33,
        "hotEvent": {"hotScore": 7000}
      }
    ]
  }]
}`

const fixtureChannel = `{
  "ret": 0,
  "idlist": [{
    "ids_hash": "def",
    "has_more": 1,
    "newslist": [
      {
        "id": "TIP2022042216544300",
        "articletype": "560",
        "title": "类型头"
      },
      {
        "id": "20260615A03TECH0",
        "articletype": "0",
        "title": "科技频道文章",
        "abstract": "科技文章摘要",
        "surl": "https://view.inews.qq.com/a/20260615A03TECH0",
        "time": "2026-06-15 09:00:00",
        "timestamp": 1749942000,
        "chlid": "news_news_tech",
        "chlname": "科技",
        "readCount": 15000,
        "hotEvent": {"hotScore": 8000}
      }
    ]
  }]
}`

const fixtureTrending = `{
  "ret": 0,
  "type": 1,
  "trace_id": "trace123",
  "topWords": {
    "fixed": [],
    "alternate": [
      {"word": "人工智能", "from": "search"},
      {"word": "世界杯", "from": "search"},
      {"word": "新冠变异", "from": "search"}
    ],
    "refreshInterval": 30,
    "displayInterval": 10
  },
  "hotlist": [],
  "tablist": []
}`

const fixtureArticleHTML = `<!DOCTYPE html>
<html>
<head><title>AI在新闻领域的最新突破</title></head>
<body>
<script type="text/javascript">
window.DATA = {"article_id":"20260615A02AHC00","article_type":"0","title":"AI在新闻领域的最新突破","desc":"人工智能正改变新闻生产方式","media":"腾讯科技","media_id":"m_1234","pubtime":"2026-06-15 10:00:00","catalog1":"科技","userAddress":"北京","contentWordsNum":1200,"article_is_pay":false,"tags":"AI,科技","interationCount":{"like":1200,"collect":300,"share":150},"originContent":{"text":"<p>人工智能正在快速发展</p><p>改变了很多行业</p>","version":"2.0"},"originAttribute":{"IMG_0":{"url":"https://inews.gtimg.com/img1.jpg","bigOrigUrl":"https://inews.gtimg.com/img1_big.jpg","width":800,"height":600,"desc":"AI示意图","islong":0},"VIDEO_0":{"vid":"v001","img":"https://pic.qq.com/cover.jpg","duration":"02:30","playurl":"https://v.qq.com/v001.mp4","title":"AI演示视频"}}};
</script>
</body>
</html>`

const fixtureVideoArticleHTML = `<!DOCTYPE html>
<html>
<body>
<script>
window.DATA = {"article_id":"20260615V09OM6","article_type":"56","title":"精彩视频","desc":"一个精彩视频","media":"腾讯视频","media_id":"m_5678","pubtime":"2026-06-15 12:00:00","catalog1":"视频","interationCount":{"like":500,"collect":50,"share":80},"originContent":{"text":"<b>精彩</b>视频内容<!--comment-->","version":"2.0"},"originAttribute":{}};
</script>
</body>
</html>`

// --- tests ---

func TestHotParsesItems(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fixtureHot))
	}))
	defer ts.Close()

	articles, err := newTestClient(ts).Hot(context.Background(), 0)
	if err != nil {
		t.Fatalf("Hot: %v", err)
	}
	if len(articles) != 2 {
		t.Fatalf("want 2 articles, got %d", len(articles))
	}
	if articles[0].ID != "20260615A02AHC00" {
		t.Errorf("first article ID = %q, want 20260615A02AHC00", articles[0].ID)
	}
	if articles[0].Rank != 1 {
		t.Errorf("first article rank = %d, want 1", articles[0].Rank)
	}
	if articles[1].Rank != 2 {
		t.Errorf("second article rank = %d, want 2", articles[1].Rank)
	}
}

func TestHotSkipsTypeHeader(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fixtureHot))
	}))
	defer ts.Close()

	articles, err := newTestClient(ts).Hot(context.Background(), 0)
	if err != nil {
		t.Fatalf("Hot: %v", err)
	}
	for _, a := range articles {
		if a.ID == "TIP2022042216544300" {
			t.Errorf("type-560 placeholder must be skipped, but got ID %q", a.ID)
		}
	}
}

func TestHotLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fixtureHot))
	}))
	defer ts.Close()

	articles, err := newTestClient(ts).Hot(context.Background(), 1)
	if err != nil {
		t.Fatalf("Hot: %v", err)
	}
	if len(articles) != 1 {
		t.Errorf("want 1 article with limit=1, got %d", len(articles))
	}
}

func TestChannelParsesItems(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fixtureChannel))
	}))
	defer ts.Close()

	articles, err := newTestClient(ts).Channel(context.Background(), "news_news_tech", 0)
	if err != nil {
		t.Fatalf("Channel: %v", err)
	}
	if len(articles) != 1 {
		t.Fatalf("want 1 article, got %d", len(articles))
	}
	if articles[0].ID != "20260615A03TECH0" {
		t.Errorf("article ID = %q, want 20260615A03TECH0", articles[0].ID)
	}
	if articles[0].HotScore != 8000 {
		t.Errorf("hot_score = %d, want 8000", articles[0].HotScore)
	}
}

func TestArticleTextDetail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fixtureArticleHTML))
	}))
	defer ts.Close()

	d, err := newTestClient(ts).Article(context.Background(), "20260615A02AHC00", false)
	if err != nil {
		t.Fatalf("Article: %v", err)
	}
	if d.ID != "20260615A02AHC00" {
		t.Errorf("ID = %q, want 20260615A02AHC00", d.ID)
	}
	if d.Title != "AI在新闻领域的最新突破" {
		t.Errorf("Title = %q", d.Title)
	}
	if d.Publisher != "腾讯科技" {
		t.Errorf("Publisher = %q, want 腾讯科技", d.Publisher)
	}
	if d.Likes != 1200 {
		t.Errorf("Likes = %d, want 1200", d.Likes)
	}
	if d.Collects != 300 {
		t.Errorf("Collects = %d, want 300", d.Collects)
	}
	if d.Body != "" {
		t.Errorf("Body should be empty when withBody=false, got %q", d.Body)
	}
	if len(d.Images) != 1 {
		t.Errorf("Images = %d, want 1", len(d.Images))
	} else {
		if d.Images[0].URL != "https://inews.gtimg.com/img1_big.jpg" {
			t.Errorf("Image URL = %q, want big URL", d.Images[0].URL)
		}
	}
	if len(d.Videos) != 1 {
		t.Errorf("Videos = %d, want 1", len(d.Videos))
	} else {
		if d.Videos[0].VID != "v001" {
			t.Errorf("Video VID = %q, want v001", d.Videos[0].VID)
		}
	}
}

func TestArticleBodyStrip(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fixtureArticleHTML))
	}))
	defer ts.Close()

	d, err := newTestClient(ts).Article(context.Background(), "20260615A02AHC00", true)
	if err != nil {
		t.Fatalf("Article: %v", err)
	}
	if d.Body == "" {
		t.Fatal("Body should be non-empty when withBody=true")
	}
	// Stripped body must not contain HTML tags
	if contains(d.Body, "<p>") || contains(d.Body, "</p>") {
		t.Errorf("Body still contains HTML tags: %q", d.Body)
	}
}

func TestArticleVideoDetail(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fixtureVideoArticleHTML))
	}))
	defer ts.Close()

	// For videos the feed ID may differ from the article_id in window.DATA — we always use the caller-supplied ID
	d, err := newTestClient(ts).Article(context.Background(), "20260615V09OM600", true)
	if err != nil {
		t.Fatalf("Article: %v", err)
	}
	if d.ID != "20260615V09OM600" {
		t.Errorf("ID = %q, want 20260615V09OM600 (caller-supplied, not the truncated window.DATA id)", d.ID)
	}
	if d.Title != "精彩视频" {
		t.Errorf("Title = %q", d.Title)
	}
	// Body should strip HTML tags and comments
	if contains(d.Body, "<b>") || contains(d.Body, "<!--") {
		t.Errorf("Body not fully stripped: %q", d.Body)
	}
}

func TestTrendingParsesWords(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fixtureTrending))
	}))
	defer ts.Close()

	words, err := newTestClient(ts).HotWords(context.Background(), 0)
	if err != nil {
		t.Fatalf("HotWords: %v", err)
	}
	if len(words) != 3 {
		t.Fatalf("want 3 words, got %d", len(words))
	}
	if words[0].Word != "人工智能" {
		t.Errorf("words[0] = %q, want 人工智能", words[0].Word)
	}
	if words[0].Rank != 1 {
		t.Errorf("words[0].Rank = %d, want 1", words[0].Rank)
	}
}

func TestTrendingLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fixtureTrending))
	}))
	defer ts.Close()

	words, err := newTestClient(ts).HotWords(context.Background(), 2)
	if err != nil {
		t.Fatalf("HotWords: %v", err)
	}
	if len(words) != 2 {
		t.Errorf("want 2 words with limit=2, got %d", len(words))
	}
}

func TestRetriesOn503(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte(fixtureHot))
	}))
	defer ts.Close()

	cfg := qq.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	cfg.Retries = 3
	client := qq.NewClient(cfg)

	articles, err := client.Hot(context.Background(), 0)
	if err != nil {
		t.Fatalf("Hot after retries: %v", err)
	}
	if len(articles) == 0 {
		t.Error("expected articles after successful retry")
	}
	if calls < 3 {
		t.Errorf("expected at least 3 calls (2 failures + 1 success), got %d", calls)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}()
}
