FROM golang:1.12.4 as builder

RUN mkdir build
COPY . /build
WORKDIR /build
RUN CC=gcc CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=auto go build -ldflags "-extldflags '-static'" -mod vendor -o runner .

FROM alpine:3.9
WORKDIR /runner
RUN apk add --no-cache --virtual .php-deps \
        php7 php7-phar php7-openssl php7-zlib php7-ctype \
        php7-mbstring php7-json php7-xml php7-dom php7-iconv \
        php7-curl php7-tokenizer php7-xmlwriter php7-simplexml \
        php7-sockets php7-bcmath php7-opcache gettext wrk curl
ADD --chown=nobody:nobody https://raw.githubusercontent.com/eficode/wait-for/828386460d138e418c31a1ebf87d9a40f5cedc32/wait-for /wait-for.sh
RUN chmod a+x /wait-for.sh
USER nobody
COPY --from=builder /build/runner .
COPY ./.rr.yml .
COPY ./php.ini /etc/php7/php.ini
EXPOSE 8080
CMD /wait-for.sh ${KAFKA_BROKER} -- ./runner serve -v -o custom.kafka.host=${KAFKA_BROKER} -o custom.kafka.topic=${KAFKA_TOPIC}