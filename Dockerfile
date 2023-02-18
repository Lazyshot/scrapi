FROM golang as builder

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go build -o scrapi .

FROM chromedp/headless-shell

RUN apt update && apt install dumb-init && rm -rf /var/lib/apt/lists/*
RUN mkdir /app
COPY --from=builder /app/scrapi /app/

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/app/scrapi", "serve"]