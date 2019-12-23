FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go install ogamego

CMD /go/bin/ogamego


EXPOSE 8080