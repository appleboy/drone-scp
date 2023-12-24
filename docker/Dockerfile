FROM alpine:3.19

ARG TARGETOS
ARG TARGETARCH

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>"

LABEL org.opencontainers.image.source=https://github.com/appleboy/drone-scp
LABEL org.opencontainers.image.description="Copy files and artifacts via SSH"
LABEL org.opencontainers.image.licenses=MIT

RUN apk add --no-cache ca-certificates && \
  rm -rf /var/cache/apk/*

RUN addgroup \
  -S -g 1000 \
  deploy && \
  adduser \
  -S -H -D \
  -h /home/deploy \
  -s /bin/sh \
  -u 1000 \
  -G deploy \
  deploy

RUN mkdir -p /home/deploy && \
  chown deploy:deploy /home/deploy

# deploy:deploy
USER 1000:1000

COPY release/${TARGETOS}/${TARGETARCH}/drone-scp /bin/

ENTRYPOINT ["/bin/drone-scp"]
