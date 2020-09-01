FROM golang:1.15.0 as builder
 
RUN useradd -u 10001 safeguard
  
WORKDIR /go/src/github.com/jw-s/safeguard
 
RUN apt-get update && \
    apt-get install ca-certificates
 
COPY . .

RUN make build

FROM scratch
VOLUME /tmp
WORKDIR /opt
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/jw-s/safeguard/bin/safeguard /opt/safeguard
USER safeguard
EXPOSE 8080
ENTRYPOINT ["/opt/safeguard", "-logtostderr"]