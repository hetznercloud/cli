FROM golang:1.21-alpine AS builder

WORKDIR /build

# Download dependencies before copying the rest of the code to take advantage of caching
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV GOCACHE=/build/cache
ENV CGO_ENABLED=0
RUN --mount=type=cache,target="/build/cache"  \
    go build  \
      -o hcloud \
      -trimpath \
      -ldflags "-s -w -X github.com/hetznercloud/cli/internal/version.versionPrerelease=" \
      ./cmd/hcloud

FROM scratch

VOLUME /config
ENV HCLOUD_CONFIG=/config/cli.toml

WORKDIR /app
COPY --from=builder /build/hcloud /app/hcloud
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/hcloud"]