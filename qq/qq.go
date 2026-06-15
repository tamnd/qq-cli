// Package qq is the library behind the qq command: HTTP client and typed data
// models for Tencent News (腾讯新闻) at news.qq.com.
//
// All metadata commands work without authentication. The Referer and User-Agent
// headers are required on every request.
//
// qq-cli is not affiliated with Tencent or QQ.
package qq

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	defaultBaseURL        = "https://i.news.qq.com"
	defaultArticleBaseURL = "https://news.qq.com"
	defaultUA             = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

// Config holds constructor parameters for Client.
type Config struct {
	BaseURL        string
	ArticleBaseURL string
	UserAgent      string
	Rate           time.Duration
	Retries        int
	Timeout        time.Duration
}

// Client fetches data from the Tencent News public API.
type Client struct {
	cfg  Config
	http *http.Client
	mu   sync.Mutex
	last time.Time
}

// NewClient returns a Client configured with cfg.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}

func (c *Client) getJSON(ctx context.Context, u string, dst any) error {
	b, err := c.get(ctx, u)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, dst); err != nil {
		return fmt.Errorf("decode %s: %w", u, err)
	}
	return nil
}

func (c *Client) get(ctx context.Context, u string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		b, retry, err := c.do(ctx, http.MethodGet, u)
		if err == nil {
			return b, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("get %s: %w", u, lastErr)
}

// articleGet fetches an article HTML page, following redirects.
func (c *Client) articleGet(ctx context.Context, articleID string) ([]byte, error) {
	u := c.cfg.ArticleBaseURL + "/rain/a/" + articleID
	return c.get(ctx, u)
}

func (c *Client) do(ctx context.Context, method, u string) ([]byte, bool, error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	req.Header.Set("Referer", "https://www.qq.com/")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	switch {
	case resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500:
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	case resp.StatusCode == http.StatusNotFound:
		return nil, false, ErrNotFound
	case resp.StatusCode == http.StatusForbidden:
		return nil, false, fmt.Errorf("http 403 forbidden")
	case resp.StatusCode != http.StatusOK:
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}

	b, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, true, err
	}
	return b, false, nil
}

func (c *Client) pace() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 500 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}
