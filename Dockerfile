FROM golang:alpine AS builder
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git bash && mkdir -p /build/paczkobot

WORKDIR /build/paczkobot

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download -json

COPY . .

RUN mkdir -p /app && CGO_ENABLED=0 go build -ldflags='-s -w -extldflags="-static"' -o /app/paczkobot

FROM scratch AS bin-unix
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/paczkobot /app/paczkobot

ENTRYPOINT ["/app/paczkobot"]
