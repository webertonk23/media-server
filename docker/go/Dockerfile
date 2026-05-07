FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache \
    git \
    gcc \
    musl-dev

RUN go install github.com/air-verse/air@latest

COPY docker/go/start.sh /start.sh

RUN chmod +x /start.sh

ENV CGO_ENABLED=1

CMD ["/start.sh"]