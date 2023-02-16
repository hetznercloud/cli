ARG GO_VERSION=1.19
ARG ALPINE_VERSION=3.16

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

# CA certificates
RUN apk add -q --update --progress --no-cache ca-certificates make curl gpg dirmngr bash sudo bat file tzdata git mandoc git-doc openssh-client zsh zsh-vcs vim libstdc++ github-cli

ENV EDITOR=vim
ENV LANG=en_US.UTF-8
ENV TERM=xterm
ENV GO111MODULE=on

# install golang tools
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.49.0
RUN go install github.com/cweill/gotests/... && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/fatih/gomodifytags@latest && \
    go install github.com/vektra/mockery/v2@latest && \
    go install github.com/golang/mock/mockgen@v1.6.0 && \
    go install golang.org/x/tools/gopls@latest && \
    go install github.com/josharian/impl@latest && \
    go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest
