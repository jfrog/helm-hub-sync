FROM golang:1.12-alpine AS builder

ADD . /root/app

RUN apk add --no-cache git

# Download modules
RUN cd /root/app && \
    GO111MODULE=on GOPROXY=https://gocenter.io go mod download

# Build microservices
RUN cd /root/app && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM gcr.io/distroless/static
COPY --from=builder /root/app/helm-hub-sync /bin

# Create user
ARG uid=1000
ARG gid=1000
RUN addgroup -g $gid helmhubsync && \
    adduser -D -u $uid -G helmhubsync helmhubsync

USER helmhubsync

CMD ["helm-hub-sync"]
