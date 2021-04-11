FROM golang:rc-buster
MAINTAINER Alex Perkins

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . . 
RUN go build -o main .

EXPOSE 3000

CMD ["/bin/bash"]
