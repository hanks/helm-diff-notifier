FROM golang:1.15.6-alpine3.12

ENV staticcheck_version=2020.2
ENV hadolint_version=v1.19.0
ENV GO111MODULE=on

WORKDIR /go/src/github.com/hanks/helm-diff-notifier

RUN apk update && \
    apk add --no-cache alpine-sdk=1.0-r0 && \
    rm -rf /var/cache/apk/*

RUN go get github.com/mattn/goveralls && \
    go get honnef.co/go/tools/cmd/staticcheck

SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

RUN wget -qO- https://github.com/hadolint/hadolint/releases/download/${hadolint_version}/hadolint-Linux-x86_64 -O /usr/local/bin/hadolint && \
    chmod +x /usr/local/bin/hadolint && \
    mkdir -p ./dist/bin

CMD ["sh"]
