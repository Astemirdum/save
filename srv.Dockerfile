FROM golang:1.18-alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY ../.. .
RUN go build -ldflags="-s -w" -o /app/srv ./server/cmd/main.go

FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app
#RUN mkdir server
#COPY --from=builder /app/srv /app/server/srv
#COPY --from=builder /build/server/cmd/server.env /app/server

COPY --from=builder /app/srv /app/srv
COPY --from=builder /build/server/cmd/server.env /app

ENTRYPOINT ["./srv"]