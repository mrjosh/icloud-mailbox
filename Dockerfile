# Golang backend environment build container
FROM golang:1.24-alpine AS build

WORKDIR /build
ADD . .

ENV GOOS linux 
ENV GOARCH amd64

RUN go mod download
RUN CGO_ENABLED=0 go build -o bin/mailbox .

# Start container
FROM alpine:3.19.0

WORKDIR /
COPY --from=build /build/bin/mailbox /usr/local/bin

ENTRYPOINT ["/usr/local/bin/mailbox"]
