FROM alpine

RUN apk add --no-cache ca-certificates

COPY kyc kyc.cfg /usr/local/bin/

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/kyc"]
