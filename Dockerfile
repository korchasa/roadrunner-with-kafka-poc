FROM golang:1.12.4-alpine as builder

# To fix go get and build with cgo
RUN apk add --no-cache --virtual .build-deps \
    bash \
    gcc \
    git \
    musl-dev

RUN mkdir build
COPY ./runner /build
WORKDIR /build

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -mod vendor -o runner .


FROM alpine:3.9
RUN apk add --no-cache --virtual .php-fpm \
        php7 php7-phar php7-openssl php7-zlib php7-ctype \
        php7-mbstring php7-json php7-xml php7-dom php7-iconv \
        php7-curl php7-tokenizer php7-xmlwriter php7-simplexml \
        php7-sockets php7-bcmath php7-opcache gettext && \
    apk add php7-rdkafka --no-cache --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing/
WORKDIR /app
COPY --from=builder /build/runner .
COPY ./php ./php
COPY ./runner/.rr.yml .
COPY ./runner/php.ini /etc/php7/php.ini
EXPOSE 8080
CMD ./runner serve -v -d