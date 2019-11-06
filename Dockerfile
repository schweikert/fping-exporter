FROM golang:alpine AS build
RUN apk add --update --no-cache ca-certificates git
ADD . /src
RUN cd /src && CGO_ENABLED=0 go build -o fping-exporter

FROM alpine:latest
RUN apk add --update --no-cache ca-certificates fping
COPY --from=build /src/fping-exporter /
EXPOSE 9605
ENTRYPOINT ["/fping-exporter", "--fping=/usr/sbin/fping"]

