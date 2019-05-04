# Start from scratch image and add in a precompiled binary
# docker build  --tag="opencoredata/ocdweb:0.1"  .
# docker run -d -p 9900:9900  opencoredata/ocdweb:0.1
FROM scratch

# Add in the static elements (could also mount these from local filesystem)
ADD cmd/services/services /

# Add our binary
CMD ["/services"]

# Document that the service listens on this port
EXPOSE 6789
