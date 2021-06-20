### --- Service --- ###
FROM alpine:latest
MAINTAINER alexperkins.crypto@gmail.com
ADD swallowtail.s.web3 /
EXPOSE 8080
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/swallowtail.s.web3"]
