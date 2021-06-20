### --- Service --- ###
FROM alpine:latest
MAINTAINER alexperkins.crypto@gmail.com
ADD swallowtail.swallowtail /
EXPOSE 8080
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/swallowtail.swallowtail"]
