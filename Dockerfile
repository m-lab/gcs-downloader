FROM golang:1.11 as build
COPY . /go/src/github.com/m-lab/gcs-downloader
ENV CGO_ENABLED 0
RUN go get -v github.com/m-lab/gcs-downloader

# Now copy the built image into the minimal base image
FROM alpine
COPY --from=build /go/bin/gcs-downloader /
# Install ca-certificates so the process can contact TLS services securely.
RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /
ENTRYPOINT ["/gcs-downloader"]
