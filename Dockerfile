# Build image
FROM golang:1.17-alpine AS builder

RUN apk update \
    && apk upgrade \
    && apk add --update \
    ca-certificates \
    gcc \
    git \
    libc-dev \
    make \
    && update-ca-certificates

WORKDIR ${GOPATH}/src/github.com/guilherme-santos/user

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN make go-install

# Final image
FROM alpine:3.11

LABEL maintainer="Guilherme Santos <xguiga@gmail.com>"

# set up nsswitch.conf for Go's "netgo" implementation
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

ENV HEALTHCHECK_VERSION 1.1.0
ENV HEALTHCHECK_URL https://github.com/gioxtech/healthcheck/releases/download/v${HEALTHCHECK_VERSION}/healthcheck-${HEALTHCHECK_VERSION}
RUN wget ${HEALTHCHECK_URL} -O /usr/bin/healthcheck && \
    chmod +x /usr/bin/healthcheck

HEALTHCHECK --start-period=5s --interval=30s --timeout=3s --retries=6 CMD healthcheck -http-addr http://localhost/health

ENV USERSVC_MYSQL_MIGRATION_DIR /etc/mysql/migrations

COPY --from=builder /go/bin/user /usr/bin/
COPY --from=builder /go/src/github.com/guilherme-santos/user/mysql/migrations /etc/mysql/migrations
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# use a non-root user
USER nobody

EXPOSE 80

CMD ["user"]
