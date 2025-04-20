FROM golang:1.24.1-alpine AS builder

ARG SERVICE

WORKDIR /go/src/github.com/htchan/BookSpider
RUN go env -w GOMODCACHE=/root/.cache/go-build

RUN --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=cache,target=/root/.cache/go-build \
    apk add gcc musl-dev libc-dev ; go mod download -x

RUN --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=internal,target=internal \
    --mount=type=bind,source=cmd,target=cmd \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -v ./cmd/$SERVICE


FROM alpine

ARG SERVICE
ENV SERVICE=$SERVICE

WORKDIR /usr/src/app

COPY --from=builder --chmod=755 /go/src/github.com/htchan/WebHistory/${SERVICE} ./app

ENTRYPOINT ["/usr/src/app/app"]