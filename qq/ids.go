package qq

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tamnd/any-cli/kit/errs"
)

var reArticleID = regexp.MustCompile(`^[0-9]{8}[AV][A-Z0-9]{7}$`)

// ParseArticleID extracts an article ID from a bare ID string or a URL.
// Accepted forms:
//   - Bare ID: "20260615A02AHC00"
//   - view.inews.qq.com URL: "https://view.inews.qq.com/a/20260615A02AHC00"
//   - news.qq.com URL: "https://news.qq.com/rain/a/20260615A02AHC00"
func ParseArticleID(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", errs.Usage("article ID or URL is required")
	}
	if reArticleID.MatchString(s) {
		return s, nil
	}
	// Extract from URL path
	for _, prefix := range []string{"/a/", "/rain/a/"} {
		if idx := strings.LastIndex(s, prefix); idx >= 0 {
			id := s[idx+len(prefix):]
			// Strip query string
			if q := strings.IndexByte(id, '?'); q >= 0 {
				id = id[:q]
			}
			if reArticleID.MatchString(id) {
				return id, nil
			}
		}
	}
	return "", errs.Usage("unrecognised article reference %q (expected ID like 20260615A02AHC00 or view.inews.qq.com URL)", s)
}

// resolveChannel returns the Channel for a user-supplied name or channel_id.
// Resolution order (case-insensitive):
//  1. Exact match on channel_id
//  2. Exact match on name_en
//  3. Exact match on name_cn
//  4. Prefix match: "tech" → "news_news_tech"
func resolveChannel(input string) (Channel, error) {
	if input == "" {
		return Channel{}, errs.Usage("channel name or ID is required")
	}
	lower := strings.ToLower(input)

	// Pass 1: exact channel_id
	for _, ch := range staticChannels {
		if strings.ToLower(ch.ChannelID) == lower {
			return ch, nil
		}
	}
	// Pass 2: exact name_en
	for _, ch := range staticChannels {
		if strings.ToLower(ch.NameEN) == lower {
			return ch, nil
		}
	}
	// Pass 3: exact name_cn
	for _, ch := range staticChannels {
		if ch.NameCN == input {
			return ch, nil
		}
	}
	// Pass 4: suffix match against channel_id (e.g. "tech" → "news_news_tech")
	for _, ch := range staticChannels {
		if strings.HasSuffix(strings.ToLower(ch.ChannelID), "_"+lower) {
			return ch, nil
		}
	}
	return Channel{}, fmt.Errorf("unknown channel %q; run 'qq channels' to see available channels", input)
}
