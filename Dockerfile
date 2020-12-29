## Builder
FROM golang:1.14-alpine3.12 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o users github.com/arnaz06/users/cmd/api

## Distribution
FROM alpine:3.12

RUN apk update && apk upgrade \
    && apk --no-cache --update add ca-certificates tzdata \
    && rm -f /var/cache/apk/*

ENV TZ Asia/Jakarta

WORKDIR /app

COPY --from=builder /app/users .


EXPOSE 7723
