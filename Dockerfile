FROM golang:alpine AS builder

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GO111MODULE=on

WORKDIR /opt/app

RUN apk --no-cache update && \
      apk --no-cache add git ca-certificates && \
      rm -rf /var/cache/apk/*

COPY . ./

RUN go build -a -o app .

FROM alpine

RUN apk --no-cache update && \
      apk --no-cache add ca-certificates && \
      rm -rf /var/cache/apk/*
      
COPY --from=builder /opt/app/app /usr/local/bin/app
RUN mkdir /usr/local/bin/log && \
      touch /usr/local/bin/log/authengine.log 

ENV LOG_FILE_LOCATION=/usr/local/bin/log/authengine.log
ENV DB_HOST=127.0.0.1
ENV DB_PORT=8669
ENV DB_USER=postgres
ENV DB_PASSWORD=Sarah4Daprinz
ENV DB_NAME=persian_black
ENV DB_SSL_MODE=disable


CMD ["/usr/local/bin/app", "--help"]