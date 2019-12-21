FROM golang:latest

RUN go get github.com/alaingilbert/ogame
RUN gop get github.com/fatih/structs

RUN .\main.go.exe