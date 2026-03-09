FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build

WORKDIR /src

ENV GOPROXY=https://proxy.golang.org,direct

RUN apk add --no-cache build-base libwebp-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=1 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
  go build -ldflags="-s -w -X github.com/flaboy/painter/internal/buildinfo.Version=${VERSION}" \
  -o /out/painter ./cmd/painter

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /out/painter /usr/local/bin/painter

EXPOSE 7013

ENTRYPOINT ["painter"]
