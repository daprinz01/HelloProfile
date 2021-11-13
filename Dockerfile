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
      rm -rf /var/cache/apk/* &&\
      mkdir /usr/local/bin/log && \
      touch /usr/local/bin/log/authengine.log
      
COPY --from=builder /opt/app/app /usr/local/bin/app
 

ENV LOG_FILE_LOCATION=/usr/local/bin/log/authengine.log \
DB_HOST=kashin.db.elephantsql.com DB_PORT=5432 DB_USER=bxywvlbh \
 DB_PASSWORD=zfIcc3jU0ajg-L7Chxx1BeNKam6FDqSp DB_NAME=bxywvlbh DB_SSL_MODE=disable \
  JWT_SECRET_KEY=QWh1aWFzdWRoaXloa1VZYmpoamFzaGQ4OTA4ODc5OHVpaEhH \
  TOKEN_LIFESPAN=5m SESSION_LIFESPAN=15m LOCK_OUT_COUNT=5 SMTP_HOST=smtp.persianblack.com \
  SMTP_PORT=25 SMTP_USER=it@persianblack.com SMTP_PASSWORD=Princess4Daprinz \
  OTP_LENGTH=6 COMMUNICATION_SERVICE_ENDPOINT=http://host.docker.internal:8084 \
  EMAIL_PATH=/api/v1/send/email OTP_VALIDITY_PERIOD=5 SMS_PATH=/api/v1/send/sms


CMD ["/usr/local/bin/app", "--help"]