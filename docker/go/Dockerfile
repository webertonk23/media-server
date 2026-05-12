FROM golang:1.25-alpine

RUN apk add --no-cache \
    ffmpeg \
    libva-utils \
    mesa-va-gallium \
    libva-intel-driver \
    git \
    gcc \
    musl-dev \
    nodejs \
    npm

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY docker/go/start.sh /start.sh

RUN chmod +x /start.sh

ENV CGO_ENABLED=1

CMD ["/start.sh"]