FROM gcr.io/distroless/base-debian10 

COPY output/acme-solver /acme-solver

USER nobody
ENTRYPOINT ["/acme-solver"]
