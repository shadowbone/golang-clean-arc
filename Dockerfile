# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
     -trimpath \
     -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
     -o /out/app ./cmd/web

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
     -trimpath -ldflags="-s -w" \
     -o /out/cleanup ./cmd/cleanup

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata wget && \
    adduser -D -u 10001 appuser

COPY --from=builder /out/app /app/app
COPY --from=builder /out/cleanup /app/cleanup

USER appuser
WORKDIR /app

EXPOSE 3000

ENTRYPOINT [ "/app/app" ]
