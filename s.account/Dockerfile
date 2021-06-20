### --- Service --- ###
FROM alpine:latest
MAINTAINER alexperkins.crypto@gmail.com
ADD swallowtail.s.account /
EXPOSE 8080
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/swallowtail.s.account"]
