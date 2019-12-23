FROM golang:latest

WORKDIR /go/src/app
COPY . .

CMD /go/bin/ogamego


EXPOSE 8080
