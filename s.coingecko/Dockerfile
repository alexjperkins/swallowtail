### --- Service --- ###
FROM alpine:latest
MAINTAINER alexperkins.crypto@gmail.com
ADD swallowtail.s.coingecko /
ADD ./config /s.coingecko/config
EXPOSE 8080
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/swallowtail.s.coingecko"]
