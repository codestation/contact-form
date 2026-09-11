FROM golang:1.27-alpine AS base

ARG GOPROXY
ENV GOPROXY=${GOPROXY}

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

FROM base AS dev

RUN set -eux; \
    addgroup -S vscode -g 1000; \
    adduser -S vscode -G vscode -u 1000; \
    mkdir -p /home/vscode/.cache/go-build /home/vscode/go/pkg/mod; \
    chown -R vscode:vscode /home/vscode

ENV HOME=/home/vscode
ENV GOCACHE=/home/vscode/.cache/go-build
ENV GOMODCACHE=/home/vscode/go/pkg/mod

USER vscode

FROM base AS builder

ARG CI_COMMIT_TAG

COPY go.mod go.sum /src/
RUN go mod download
COPY . /src/

RUN set -ex; \
    CGO_ENABLED=0 go build -o release/contact-form \
    -trimpath \
    -tags viper_yaml3 \
    -ldflags "-w -s \
    -X megpoid.dev/go/contact-form/version.Tag=${CI_COMMIT_TAG}"

FROM alpine:3.24 AS prod
LABEL maintainer="codestation <codestation@megpoid.dev>"

RUN apk add --no-cache ca-certificates tzdata

RUN set -eux; \
    addgroup -S runner -g 1000; \
    adduser -S runner -G runner -u 1000

COPY --from=builder /src/release/contact-form /usr/local/bin/contact-form

USER runner

EXPOSE 8000

ENTRYPOINT ["/usr/local/bin/contact-form"]
CMD ["serve"]
