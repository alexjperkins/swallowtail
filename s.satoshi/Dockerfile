FROM golang:rc-buster
MAINTAINER alexperkins.crypto@gmail.com

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . . 
RUN go build -o main .

EXPOSE 8080

CMD ["/bin/bash"]
