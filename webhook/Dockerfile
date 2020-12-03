FROM alpine:3.12.1

RUN apk add --no-cache ca-certificates

COPY output/webhook /usr/local/bin/webhook
USER nobody
ENTRYPOINT ["webhook"]
