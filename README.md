# SCRAPI - Scraper API

This project is meant to provide a simple API that can be intergrated
into workflow automation software to scrape information for a website.

The primary goal is in converting pages of listed items to structure JSON
data.

## Getting Started

There are varied methods of deploying the API.

### Kubernetes

TODO: Helm Chart
TODO: Basic deploy manifests

### Docker

A docker image is built, which includes an installation
of chromium and all other dependencies.

```sh
$ docker run -p 8000:8000 lazyshot/scrapi
```

### Binary

TODO: Release process
```sh
$ wget ...
$ ./scrapi serve
```

### Source Code

Being a go application, you would need to have a the go compiler
installed.

```sh
$ go clone https://github.com/lazyshot/scrapi.git
$ cd scrapi
$ go build .
$ ./scrapi serve
```

## API Docs

Swagger documentation is generated using swaggo, which is then
served at `/swagger/`. This is a helpful guide to exactly how to use the API.

Swagger UI: https://localhost:8000/swagger/index.html