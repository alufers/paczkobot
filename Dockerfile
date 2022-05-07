FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git bash && mkdir -p /build/paczkobot

WORKDIR /build/paczkobot

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download -json

COPY . .

RUN mkdir -p /app && CGO_ENABLED=0 GOOS=${TARGETPLATFORM%%/*} GOARCH=${TARGETPLATFORM##*/} \
    go build -ldflags='-s -w -extldflags="-static"' -o /app/paczkobot

# RUN echo "Running on architecture: $(uname -m), BUILDPLATFORM=$BUILDPLATFORM, TARGETPLATFORM=$TARGETPLATFORM" && exit 1

FROM scratch AS bin-unix
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/paczkobot /app/paczkobot

LABEL org.opencontainers.image.description A docker image for the paczkobot telegram bot.

ENTRYPOINT ["/app/paczkobot"]
