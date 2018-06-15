FROM alpine:3.6
RUN apk add --no-cache ca-certificates
ADD bin/safeguard /opt/safeguard
EXPOSE 443
ENTRYPOINT ["/opt/safeguard"]
VOLUME ["/tmp"]