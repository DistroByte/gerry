# Gerry

Platform agnostic chatbot.

[![Go report](https://goreportcard.com/badge/git.dbyte.xyz/distro/gerry)](http://goreportcard.com/report/git.dbyte.xyz/distro/gerry)

## Getting Started

### Generate Config

```bash
$ make config
```

This will output a `config.yaml` file in the root of the project. Edit this file to your liking.

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

Images are available on [GitHub Container Registry](https://github.com/DistroByte/gerry/pkgs/container/gerry) as `ghcr.io/distrobyte/gerry`.

The `latest` tag is the most recent release. Use a SHA for a specific release.
