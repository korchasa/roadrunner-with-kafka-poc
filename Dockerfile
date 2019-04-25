FROM golang:1.12.4-alpine as builder

# To fix go get and build with cgo
RUN apk add --no-cache --virtual .build-deps \
    bash \
    gcc \
    git \
    musl-dev

RUN mkdir build
COPY . /build
WORKDIR /build

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -mod vendor -o webserver .
RUN adduser -S -D -H -h /build webserver
USER webserver

FROM scratch
WORKDIR /app
COPY --from=builder /build/webserver .
COPY ./php ./php
EXPOSE 8080
CMD ["./webserver"]