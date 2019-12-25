FROM golang:latest

USER nobody

RUN mkdir -p /go/src/github.com/lordbentner/ogamego
WORKDIR /go/src/github.com/lordbentner/ogamego

COPY . /go/src/github.com/lordbentner/ogamego
#RUN go build

CMD ["./main.go"]
EXPOSE 8080
