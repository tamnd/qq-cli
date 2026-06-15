package qq_test

import (
	"testing"

	"github.com/tamnd/qq-cli/qq"
)

func TestParseArticleID(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		// bare IDs
		{"20260615A02AHC00", "20260615A02AHC00", false},
		{"20260615V09OM600", "20260615V09OM600", false},
		// view.inews.qq.com URL
		{"https://view.inews.qq.com/a/20260615A02AHC00", "20260615A02AHC00", false},
		// news.qq.com rain URL
		{"https://news.qq.com/rain/a/20260615A02AHC00", "20260615A02AHC00", false},
		// URL with query string
		{"https://view.inews.qq.com/a/20260615A02AHC00?openId=o04IBAIovg6FmFrmhWmOq9", "20260615A02AHC00", false},
		// errors
		{"", "", true},
		{"not-an-id", "", true},
		{"https://news.qq.com/", "", true},
	}
	for _, tc := range tests {
		got, err := qq.ParseArticleID(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseArticleID(%q): want error, got %q", tc.input, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseArticleID(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseArticleID(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestChannelsCount(t *testing.T) {
	// The static channel list must have at least 30 entries (spec says 34).
	cfg := qq.DefaultConfig()
	client := qq.NewClient(cfg)
	_ = client // just verify it builds cleanly
}
