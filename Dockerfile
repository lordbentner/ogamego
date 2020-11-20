# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.13 as builder

# Copy local code to the container image.
WORKDIR /ogamego
COPY . .

RUN go get github.com/alaingilbert/ogame
RUN go get github.com/fatih/structs
RUN go get github.com/go-macaron/binding
RUN go get github.com/kardianos/service
RUN go get github.com/stretchr/stew/slice
RUN go get gopkg.in/macaron.v1
RUN go get github.com/Masterminds/sprig
# Build the command inside the container.
RUN CGO_ENABLED=0 GOOS=windows go build .

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
#COPY --from=builder /ogamebot/main.exe /main.exe

# Run the web service on container startup.
CMD ["/main.exe"]