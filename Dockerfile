FROM docker.io/library/alpine:3.16 as runtime

ENTRYPOINT ["stiebeleltron-exporter"]

RUN \
    apk add --no-cache curl bash

COPY stiebeleltron-exporter /usr/bin/
USER 1000:0
