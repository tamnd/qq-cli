package qq

import (
	"context"
	"fmt"
	neturl "net/url"
)

// Hot fetches the global hot ranking list.
// Returns up to limit articles (0 = all, max 50). Skips the type-560 header item.
func (c *Client) Hot(ctx context.Context, limit int) ([]Article, error) {
	u := fmt.Sprintf("%s/gw/event/hot_ranking_list?page_size=50", c.cfg.BaseURL)
	return c.fetchFeed(ctx, u, limit)
}

// Channel fetches the top articles for a specific channel by channel_id.
// Returns up to limit articles (0 = all, max 10). Skips the type-560 header item.
func (c *Client) Channel(ctx context.Context, channelID string, limit int) ([]Article, error) {
	u := fmt.Sprintf("%s/gw/event/pc_hot_ranking_list?channelId=%s&page_size=10",
		c.cfg.BaseURL, neturl.QueryEscape(channelID))
	return c.fetchFeed(ctx, u, limit)
}

func (c *Client) fetchFeed(ctx context.Context, u string, limit int) ([]Article, error) {
	var resp wireHotResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	if len(resp.IDList) == 0 {
		return nil, nil
	}
	items := resp.IDList[0].NewsList
	out := make([]Article, 0, len(items))
	rank := 1
	for _, item := range items {
		if item.ArticleType == "560" {
			continue // skip pinned header
		}
		out = append(out, articleFrom(item, rank))
		rank++
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

// HotWords fetches the trending search words.
// Returns up to limit words (0 = all, max 60) from topWords.alternate, ranked 1-60.
func (c *Client) HotWords(ctx context.Context, limit int) ([]HotWord, error) {
	u := fmt.Sprintf("%s/gw/pc_search/hotWord", c.cfg.BaseURL)
	var resp wireHotWordResp
	if err := c.getJSON(ctx, u, &resp); err != nil {
		return nil, err
	}
	words := resp.TopWords.Alternate
	if limit > 0 && limit < len(words) {
		words = words[:limit]
	}
	out := make([]HotWord, 0, len(words))
	for i, w := range words {
		if w.Word == "" {
			continue
		}
		out = append(out, hotWordFrom(w, i+1))
	}
	return out, nil
}

// Article fetches full detail for one article by its ID.
// withBody controls whether the body text (stripped HTML) is included.
func (c *Client) Article(ctx context.Context, articleID string, withBody bool) (ArticleDetail, error) {
	html, err := c.articleGet(ctx, articleID)
	if err != nil {
		return ArticleDetail{}, err
	}
	w, err := parseWindowData(string(html))
	if err != nil {
		return ArticleDetail{}, fmt.Errorf("article %s: %w", articleID, err)
	}
	return articleDetailFrom(*w, articleID, withBody), nil
}
