# Gerry

Platform agnostic chatbot.

[![Go report](https://goreportcard.com/badge/github.com/distrobyte/gerry)](http://goreportcard.com/report/github.com/distrobyte/gerry)

## Getting Started

### Run

```bash
$ make run
```

## Docker

### Run

```bash
$ docker run --rm -v "$(pwd)/config.yaml:/app/config.yaml" ghcr.io/distrobyte/gerry:latest
```

### Images

Images are available on [GitHub Container Registry](https://github.com/distrobyte/gerry/pkgs/container/gerry) as `ghcr.io/distrobyte/gerry`.

The `latest` tag is the most recent release. Use a SHA for a specific release.

#TODO: add kek counter
