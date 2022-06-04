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
      mkdir /var/local/bin/log && \
      touch /var/local/bin/log/helloprofile.log
      
COPY --from=builder /opt/app/app /usr/local/bin/app
COPY --from=builder /opt/app/gcp-storage-config.json /usr/local/bin/gcp-storage-config.json
 

ENV LOG_FILE_LOCATION=/var/local/bin/log/helloprofile.log \
DB_HOST=drona.db.elephantsql.com DB_PORT=5432 DB_USER=iuyegkoq \
 DB_PASSWORD=Zm032Nh7TJO_A_vifLUk8gX1R49YPEMe DB_NAME=iuyegkoq DB_SSL_MODE=disable \
  JWT_SECRET_KEY=QWh1aWFzdWRoaXloa1VZYmpoamFzaGQ4OTA4ODc5OHVpaEhH \
  TOKEN_LIFESPAN=8760h SESSION_LIFESPAN=17520h LOCK_OUT_COUNT=5 SMTP_HOST=smtp.persianblack.com \
  SMTP_PORT=25 SMTP_USER=it@persianblack.com SMTP_PASSWORD=Princess4Daprinz \
  OTP_LENGTH=6 COMMUNICATION_SERVICE_ENDPOINT=http://host.docker.internal:8084 \
  EMAIL_PATH=/api/v1/send/email OTP_VALIDITY_PERIOD=60 SMS_PATH=/api/v1/send/sms \
  GOOGLE_CLIENT_ID_WEB=954287935322-knnckb397mk1uo98av0km4cbgilrm4ke.apps.googleusercontent.com \
  GOOGLE_CLIENT_ID_IOS=954287935322-o578q1g3qa042rs91p50rb92t75on2e4.apps.googleusercontent.com \
  GOOGLE_CLIENT_ID_ANDROID=954287935322-7nmsff5amchso29k2csvpsqrjqgkjv1l.apps.googleusercontent.com \
  GOOGLE_APPLICATION_CREDENTIALS=/usr/local/bin/gcp-storage-config.json GCP_BUCKET_NAME=helloprofile-test GCP_BUCKET_PROJECT_ID=bionic-region-347614 \
  GCP_UPLOAD_PATH=test/ GCP_PUBLIC_HOST=http://storage.googleapis.com

CMD ["/usr/local/bin/app", "--help"]