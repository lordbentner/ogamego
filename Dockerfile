FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get github.com/alaingilbert/ogame
RUN go get github.com/fatih/structs

RUN .\main.go.exe

EXPOSE 8080