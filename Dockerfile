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
CMD ["helm-hub-sync"]
