FROM golang:rc-buster
MAINTAINER Alex Perkins

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o main .

WORKDIR /svc

RUN cp /build/main .

EXPOSE 3000

CMROM
