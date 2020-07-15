FROM alpine:3.12

RUN apk add --no-cache ca-certificates

COPY output/webhook /usr/local/bin/webhook

ENTRYPOINT ["webhook"]
