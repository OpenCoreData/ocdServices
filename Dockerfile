# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
# docker build  --tag="opencoredata/ocdservices:0.1"  .
# docker run -d -p 9900:9900  opencoredata/ocdservices:0.1
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/opencoredata.org/ocdServices

# Create a non-root user to run as

RUN groupadd -r gorunner -g 433 && \
mkdir /home/gorunner && \
useradd -u 431 -r -g gorunner -d /home/gorunner -s /sbin/nologin -c "User to run go apps on high ports" gorunner && \
chown -R gorunner:gorunner  /home/gorunner && \
chown -R gorunner:gorunner /go/src/opencoredata.org/ocdServices


# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get code.google.com/p/gorilla/mux
#RUN go install opencoredata.org/ocdServices

# set user
USER gorunner

# Move to a workign directory for running codices so it can see it's static files
# future version should take this as a param so static content can be anywhere
WORKDIR /go/src/opencoredata.org/ocdServices


# Run the command by default when the container starts.
ENTRYPOINT go run main.go

# Document that the service listens on this port
EXPOSE 6789