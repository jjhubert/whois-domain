FROM golang:1.19 as builder
ENV GO111MODULE=on GOARCH=amd64 GOOS=linux CGO_ENABLED=0 GOPROXY=https://goproxy.io
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o whois-domain ./cmd

FROM alpine:latest
WORKDIR /app
RUN mkdir config
COPY --from=builder /app/whois-domain whois-domain
COPY --from=builder /app/config/config.yaml config/config.yaml
ENTRYPOINT  ["/app/whois-domain"]
