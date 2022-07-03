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
RUN go build -ldflags="-s -w" -o /app/client ./client/cmd/main.go

FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

#RUN mkdir client
#COPY --from=builder /app/client /app/client/client
#COPY --from=builder /build/client/cmd/client.env /app/client

COPY --from=builder /app/client /app/client
COPY --from=builder /build/client/cmd/client.env /app


ENTRYPOINT ["./client"]
