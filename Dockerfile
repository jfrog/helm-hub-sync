FROM golang:alpine as builder

ADD . /root/app

RUN apk add --no-cache curl wget git alpine-sdk \
    && cd /root/app \
    && go get ./... \
    && CGO_ENABLED=0 go build

FROM alpine:3.9
COPY --from=builder /root/app/helm-hub-sync /bin
RUN apk add --no-cache curl wget git ca-certificates
CMD ["helm-hub-sync"]