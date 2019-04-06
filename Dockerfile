# We use a so called two stage build.
# Basically this means we build our go binary in one image
# which has the go compiler and all the required tools and libraries.
# However, since we do not need those in our production image,
# we copy the binary created into a basic alpine image
# resulting in a much smaller image for production.

# We define our image for the build environment...
FROM golang:1.11-alpine3.8 as build

# ...and copy our project directory tp the according place.
COPY . "/go/src/github.com/mwmahlberg/so-postgres-compose"
# We set our work directory...
WORKDIR /go/src/github.com/mwmahlberg/so-postgres-compose
# ...and add git, which - forever reason, is not included into the golang image.
RUN apk add git

# We set up our dependency management and make sure all dependencies outside
# the standard library are installed.
RUN set -x && \
    go get github.com/golang/dep/cmd/dep && \
    dep ensure -v
# Finally, we build our binary and name it accordingly    
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /apiserver

# Now for the second stage, we use a basic alpine image...
FROM alpine:3.8
# ... update it...
RUN apk --no-cache upgrade
# .. and take the binary from the image we build above.
# Note the location: [/usr]{/bin:/sbin} are reserved for
# the OS's package manager. Binaries added to an OS by
# the administrator which are not part of the OS's package
# mangement system should always go below /usr/local/{bin,sbin}
COPY --from=build /apiserver /usr/local/bin/apiserver
# Last but not least, we tell docker which binary to execute.
ENTRYPOINT ["/usr/local/bin/apiserver"]