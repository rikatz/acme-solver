FROM alpine:3.12

RUN apk add --no-cache ca-certificates

COPY output/acme-solver /usr/local/bin/acme-solver

ENTRYPOINT ["acme-solver"]
