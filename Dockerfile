FROM golang:latest

USER nobody

RUN mkdir -p /go/src/github.com/lordbentner/ogamego
WORKDIR /go/src/github.com/lordbentner/ogamego

COPY . /go/src/github.com/lordbentner/ogamego
#RUN go build

CMD ["sudo","chmod","+x","./main.go"]
CMD ["./main.go"]
EXPOSE 8080
