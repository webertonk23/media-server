FROM golang:1.25-alpine

RUN apk add --no-cache \
    ffmpeg \
    git \
    gcc \
    musl-dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY docker/go/start.sh /start.sh

RUN chmod +x /start.sh

ENV CGO_ENABLED=1

CMD ["/start.sh"]