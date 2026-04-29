# SlugGo

Simple URL Shortener built with Golang

## Commands

This project uses a `Makefile` as the closest equivalent to `npm scripts` in Go.

```bash
make run-local
make run-docker
make update-swagger
```

These targets simply call the scripts in the `scripts/` directory, so the shell scripts remain the source of truth.
