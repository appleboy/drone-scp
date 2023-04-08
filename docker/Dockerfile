FROM alpine:3.17

ARG TARGETOS
ARG TARGETARCH

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="SCP Plugin" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

LABEL org.opencontainers.image.source=https://github.com/appleboy/drone-scp
LABEL org.opencontainers.image.description="Copy files and artifacts via SSH"
LABEL org.opencontainers.image.licenses=MIT

RUN apk add --no-cache ca-certificates=20220614-r4 && \
  rm -rf /var/cache/apk/*

COPY release/${TARGETOS}/${TARGETARCH}/drone-scp /bin/

ENTRYPOINT ["/bin/drone-scp"]
