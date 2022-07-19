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
      touch /usr/local/bin/log/helloprofile.log
      
COPY --from=builder /opt/app/app /usr/local/bin/app
 

ENV LOG_FILE_LOCATION=/usr/local/bin/log/helloprofile.log \
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
  FIREBASE_SERVICE_ACCOUNT_CREDENTIAL_PATH=/usr/local/bin/serviceAccountKey.json
  PROFILE_TEMPLATES=c49450f6-e1bb-43de-9b4f-519fabb0eb6e,b1fefdcb-caca-4b16-baf6-b5bf21a54356,2234ef30-b756-403c-958e-b56e11038d1a,1fea9c1e-8911-4a99-a522-1310fb99e74c \
  SPACES_KEY=key SPACES_SECRET=secret DIGITAL_OCEAN_SPACES_ENDPOINT=https://fra1.digitaloceanspaces.com DIGITAL_OCEAN_SPACE_NAME=helloprofile \
  DIGITAL_OCEAN_SPACES_CDN=https://helloprofile.fra1.cdn.digitaloceanspaces.com
CMD ["/usr/local/bin/app", "--help"]