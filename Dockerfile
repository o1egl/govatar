FROM alpine:latest as alpine
RUN apk --update add ca-certificates
RUN apk --update add mailcap

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=alpine /etc/mime.types /etc/mime.types

VOLUME /data

COPY govatar /govatar

ENTRYPOINT [ "/govatar" ]