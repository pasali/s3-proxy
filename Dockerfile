FROM golang:1.12-alpine
ADD s3-proxy /usr/sbin/
EXPOSE 8080
ENTRYPOINT ["s3-proxy"]
