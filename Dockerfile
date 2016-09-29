# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
# docker build  --tag="opencoredata/ocdservices:0.1"  .
# docker run -d -p 6789:6789  opencoredata/ocdservices:0.1
FROM golang:1.5

# Copy the local package files to the container's workspace.
ADD . /go/src/opencoredata.org/ocdServices
#ADD ./ocdCommons /go/src/opencoredata.org/ocdCommons
ADD ./buildDependencies /oracle
#ADD https://codeload.github.com/OpenCoreData/ocdCommons/zip/master  /

# Uncompress ocdCommons
# code here to uncompress and move the commons package

# Create a non-root user to run as
RUN groupadd -r gorunner -g 433 && \
mkdir /home/gorunner && \
useradd -u 431 -r -g gorunner -d /home/gorunner -s /sbin/nologin -c "User to run go apps on high ports" gorunner && \
chown -R gorunner:gorunner  /home/gorunner && \
chown -R gorunner:gorunner /go/src/opencoredata.org/ocdServices

# add in some exports I suspect we need to build with
RUN apt-get update && apt-get install -y pkg-config libaio-dev libaio1
RUN export PKG_CONFIG_PATH=/oracle
ENV PKG_CONFIG_PATH /oracle

RUN export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/oracle/instantclient_12_1
ENV LD_LIBRARY_PATH $LD_LIBRARY_PATH:/oracle/instantclient_12_1

RUN export CGO_CFLAGS=-I/oracle/instantclient_12_1/sdk/include
RUN export CGO_LDFLAGS="-L/oracle/instantclient_12_1 -lclntsh"

ENV CGO_CFLAGS -I/oracle/instantclient_12_1/sdk/include
ENV CGO_LDFLAGS "-L/oracle/instantclient_12_1 -lclntsh"

#  mac only line export DYLD_LIBRARY_PATH=/oracle/instantclient_11_2:$DYLD_LIBRARY_PATH

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get gopkg.in/rana/ora.v3
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2
RUN go get github.com/emicklei/go-restful
RUN go get github.com/lib/pq
RUN go get github.com/knakk/sparql
RUN go get github.com/mb0/wkt
RUN go get github.com/kpawlik/geojson
RUN go get github.com/chris-ramon/graphql
RUN go get github.com/gocarina/gocsv
RUN go get github.com/jmoiron/sqlx
RUN go get github.com/kisielk/sqlstruct
#RUN go install opencoredata.org/ocdServices

# set user
USER gorunner

# Move to a workign directory for running codices so it can see it's static files
# future version should take this as a param so static content can be anywhere
WORKDIR /go/src/opencoredata.org/ocdServices


# try and go build and run from the static version
RUN go build

# Run the command by default when the container starts.
#ENTRYPOINT go run main.go
CMD ["/go/src/opencoredata.org/ocdServices/ocdServices"]

# Document that the service listens on this port
# container needs to talk to database container
EXPOSE 6789
