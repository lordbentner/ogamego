FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN main.exe

EXPOSE 8080