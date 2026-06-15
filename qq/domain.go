package qq

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

func init() { kit.Register(Domain{}) }

// Domain is the Tencent News driver for the any-cli/kit framework.
type Domain struct{}

// Info describes the scheme and hostnames this domain handles.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme:   "qq",
		Hosts:    []string{"qq.com", "news.qq.com", "i.news.qq.com", "view.inews.qq.com"},
		Identity: BaseIdentity(),
	}
}

// BaseIdentity is the help and version identity for the qq binary.
func BaseIdentity() kit.Identity {
	return kit.Identity{
		Binary: "qq",
		Short:  "A command-line for Tencent News / QQ.com (腾讯新闻).",
		Long: `qq reads Tencent News (腾讯新闻) data and prints clean, pipeable records.

Browse the hot ranking, explore 34 channels, fetch full article text, and
view trending search words. All commands work without an account.

Records come out as table, list, markdown, JSON, JSONL, CSV, TSV, url, or raw.

qq is an independent tool and is not affiliated with Tencent or QQ.`,
		Site: "https://news.qq.com",
		Repo: "https://github.com/tamnd/qq-cli",
	}
}

// Defaults seeds the framework baseline from qq defaults.
func Defaults(c *kit.Config) {
	d := DefaultConfig()
	c.Rate = d.Rate
	c.Timeout = d.Timeout
	c.Retries = d.Retries
	c.UserAgent = d.UserAgent
}

// DefaultConfig returns a Config with sensible defaults for the Tencent News API.
func DefaultConfig() Config {
	return Config{
		BaseURL:        defaultBaseURL,
		ArticleBaseURL: defaultArticleBaseURL,
		UserAgent:      defaultUA,
		Rate:           300 * time.Millisecond,
		Timeout:        30 * time.Second,
		Retries:        3,
	}
}

// Register installs the client factory and all operations onto app.
func (Domain) Register(app *kit.App) {
	app.SetClient(newClient)
	registerOps(app)
}

// Register is a convenience so callers don't need to name the zero-value Domain.
func Register(app *kit.App) { Domain{}.Register(app) }

// Session is the per-run client that kit injects into every operation.
type Session struct {
	Client *Client
	Quiet  bool
}

// Progressf prints a one-line progress note to stderr unless quiet.
func (s *Session) Progressf(format string, args ...any) {
	if s == nil || s.Quiet {
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func newClient(_ context.Context, c kit.Config) (any, error) {
	cfg := DefaultConfig()
	if c.UserAgent != "" {
		cfg.UserAgent = c.UserAgent
	}
	if c.Rate > 0 {
		cfg.Rate = c.Rate
	}
	if c.Timeout > 0 {
		cfg.Timeout = c.Timeout
	}
	if c.Retries > 0 {
		cfg.Retries = c.Retries
	}
	return &Session{Client: NewClient(cfg), Quiet: c.Quiet}, nil
}

// MapErr converts a library error into the kit error kind with the right exit code.
func MapErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, ErrNotFound):
		return errs.NotFound("%s", err.Error())
	default:
		return err
	}
}

// Classify turns an input URL into the canonical (uriType, id).
func (Domain) Classify(input string) (uriType, id string, err error) {
	if input == "" {
		return "", "", errs.Usage("empty qq reference")
	}
	if strings.Contains(input, "/a/") || strings.Contains(input, "/rain/a/") {
		if aid, e := ParseArticleID(input); e == nil {
			return "article", aid, nil
		}
	}
	return "", "", errs.Usage("unrecognised qq reference: %q", input)
}

// Locate returns the live https URL for a (uriType, id).
func (Domain) Locate(uriType, id string) (string, error) {
	switch uriType {
	case "article":
		return "https://view.inews.qq.com/a/" + id, nil
	default:
		return "", errs.Usage("qq has no resource type %q", uriType)
	}
}
