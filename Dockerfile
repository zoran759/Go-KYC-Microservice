FROM alpine

RUN apk add --no-cache ca-certificates

COPY kyc kyc.cfg /usr/local/bin/

WORKDIR /usr/local/bin

EXPOSE 8080

ENTRYPOINT ["kyc"]
