FROM golang:buster
MAINTAINER Alex Perkins

ENV GO111MODULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

WORKDIR /app
COPY . /app

RUN go build -o main .
CMD ["/app/main"]
