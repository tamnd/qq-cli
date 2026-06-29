# qq-cli

A fast, friendly command line for Tencent News / QQ.com (腾讯新闻).

`qq` reads public Tencent News data and prints clean, pipeable records. All commands work without an account or API key.

> **Not affiliated with Tencent or QQ.** This is an independent open-source tool.

## Install

```sh
# Homebrew (macOS / Linux)
brew install tamnd/tap/qq

# Go install
go install github.com/tamnd/qq-cli/cmd/qq@latest

# Docker
docker run --rm ghcr.io/tamnd/qq:latest hot
```

Or download a pre-built binary from [Releases](https://github.com/tamnd/qq-cli/releases).

## Commands

### `qq hot` — global hot ranking

Fetch up to 50 articles from the global hot ranking list.

```sh
qq hot
qq hot -n 10
qq hot -o table
qq hot -o url
```

### `qq channel <name>` — channel feed

Fetch top articles for one of the 34 QQ News channels. Pass the English name, Chinese name, or channel ID.

```sh
qq channel tech
qq channel 科技
qq channel news_news_sports -o table
qq channel finance -n 5
```

### `qq channels` — list all channels

Print the full channel directory (34 entries).

```sh
qq channels
qq channels -o table
qq channels -o json | jq '.[] | select(.has_feed)'
```

### `qq article <id-or-url>` — article detail

Fetch full metadata for one article. Use `--body` to include the stripped body text.

```sh
qq article 20260615A02AHC00
qq article https://view.inews.qq.com/a/20260615A02AHC00
qq article 20260615A02AHC00 --body
qq article 20260615A02AHC00 -o json | jq .images
```

### `qq trending` — trending search words

Fetch up to 60 rotating trending search words.

```sh
qq trending
qq trending -n 20
qq trending -o table
```

## Output formats

Every command accepts `-o`:

| Flag         | Description                        |
|--------------|------------------------------------|
| `-o table`   | Aligned columns (default)          |
| `-o list`    | Key: value pairs                   |
| `-o json`    | Pretty-printed JSON array          |
| `-o jsonl`   | One JSON object per line (default) |
| `-o csv`     | CSV with header row                |
| `-o tsv`     | TSV with header row                |
| `-o url`     | One URL per line                   |
| `-o raw`     | Raw API response                   |

Use `-n N` to limit output to N records.

## License

Apache-2.0. See [LICENSE](LICENSE).
