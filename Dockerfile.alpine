FROM alpine:3.6

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="Drone SCP" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

RUN apk add -U --no-cache ca-certificates && \
  rm -rf /var/cache/apk/*

ADD release/linux/amd64/drone-scp /bin/

ENTRYPOINT ["/bin/drone-scp"]
