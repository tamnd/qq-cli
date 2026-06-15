package qq

import (
	"context"

	"github.com/tamnd/any-cli/kit"
)

func registerOps(app *kit.App) {
	registerHot(app)
	registerChannel(app)
	registerChannels(app)
	registerArticle(app)
	registerTrending(app)
}

// hot -----------------------------------------------------------------------

type hotIn struct {
	Session *Session `kit:"inject"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerHot(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "hot",
		Group:   "read",
		List:    true,
		Summary: "Global hot ranking list (up to 50 articles)",
	}, func(ctx context.Context, in hotIn, emit func(Article) error) error {
		in.Session.Quiet = in.Quiet
		in.Session.Progressf("fetching hot ranking")
		articles, err := in.Session.Client.Hot(ctx, 0)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(articles, emit)
	})
}

// channel -------------------------------------------------------------------

type channelIn struct {
	Session *Session `kit:"inject"`
	Name    string   `kit:"arg"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerChannel(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "channel",
		Group:   "read",
		List:    true,
		Summary: "Top articles for a specific channel",
		Args:    []kit.Arg{{Name: "name", Help: "channel name (English/Chinese) or channel_id; run 'qq channels' to list all"}},
	}, func(ctx context.Context, in channelIn, emit func(Article) error) error {
		in.Session.Quiet = in.Quiet
		ch, err := resolveChannel(in.Name)
		if err != nil {
			return err
		}
		if !ch.HasFeed {
			return MapErr(ErrNotFound)
		}
		in.Session.Progressf("fetching channel %s (%s)", ch.ChannelID, ch.NameCN)
		articles, err := in.Session.Client.Channel(ctx, ch.ChannelID, 0)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(articles, emit)
	})
}

// channels ------------------------------------------------------------------

type channelsIn struct {
	Session *Session `kit:"inject"`
}

func registerChannels(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "channels",
		Group:   "read",
		List:    true,
		Summary: "List all 34 QQ News channels",
	}, func(ctx context.Context, in channelsIn, emit func(Channel) error) error {
		return emitAll(staticChannels, emit)
	})
}

// article -------------------------------------------------------------------

type articleIn struct {
	Session *Session `kit:"inject"`
	Ref     string   `kit:"arg"`
	Body    bool     `kit:"flag"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerArticle(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "article",
		Group:   "read",
		Single:  true,
		Summary: "Fetch full article detail",
		Args:    []kit.Arg{{Name: "id", Help: "article ID (20260615A02AHC00) or view.inews.qq.com URL"}},
	}, func(ctx context.Context, in articleIn, emit func(ArticleDetail) error) error {
		in.Session.Quiet = in.Quiet
		id, err := ParseArticleID(in.Ref)
		if err != nil {
			return err
		}
		in.Session.Progressf("fetching article %s", id)
		detail, err := in.Session.Client.Article(ctx, id, in.Body)
		if err != nil {
			return MapErr(err)
		}
		return emit(detail)
	})
}

// trending ------------------------------------------------------------------

type trendingIn struct {
	Session *Session `kit:"inject"`
	Quiet   bool     `kit:"flag,inherit"`
}

func registerTrending(app *kit.App) {
	kit.Handle(app, kit.OpMeta{
		Name:    "trending",
		Group:   "read",
		List:    true,
		Summary: "60 rotating trending search words",
	}, func(ctx context.Context, in trendingIn, emit func(HotWord) error) error {
		in.Session.Quiet = in.Quiet
		in.Session.Progressf("fetching trending words")
		words, err := in.Session.Client.HotWords(ctx, 0)
		if err != nil {
			return MapErr(err)
		}
		return emitAll(words, emit)
	})
}

// helpers -------------------------------------------------------------------

func emitAll[T any](items []T, emit func(T) error) error {
	for _, item := range items {
		if err := emit(item); err != nil {
			return err
		}
	}
	return nil
}
