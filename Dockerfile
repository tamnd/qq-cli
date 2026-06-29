# GoReleaser stages the prebuilt binary under $TARGETPLATFORM in the build
# context; COPY selects the right one via the automatic TARGETPLATFORM build arg.
FROM alpine:3.21

ARG TARGETPLATFORM

RUN apk add --no-cache ca-certificates tzdata \
 && adduser -D -H -u 10001 qq \
 && mkdir -p /data \
 && chown qq:qq /data

COPY $TARGETPLATFORM/qq /usr/bin/qq

USER qq
WORKDIR /data

ENTRYPOINT ["/usr/bin/qq"]
