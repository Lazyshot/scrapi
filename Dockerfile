FROM golang as builder

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go build -o scrapi .

FROM ubuntu

RUN apt update && \
    apt install -y wget && \
    wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    apt install -y ./google-chrome-stable_current_amd64.deb && \
    apt install -y \
        dumb-init xvfb ca-certificates \
        && \
    rm -rf /var/lib/apt/lists/*
RUN mkdir /app
COPY --from=builder /app/scrapi /app/

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/app/scrapi", "serve"]