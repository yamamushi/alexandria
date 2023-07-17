# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/yamamushi/alexandria

# Create our shared volume
RUN mkdir /data

# Build the EQB command inside the container.
RUN cd /go/src/github.com/yamamushi/alexandria && go get -v ./... && go build -v ./... && go install

# Run the EQB command by default when the container starts.
WORKDIR /data
ENTRYPOINT /go/bin/alexandria

# Set the working directory to /data/
VOLUME /data
WORKDIR /data
