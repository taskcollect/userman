# PRODUCTION DOCKERFILE

# --- Build Stage

FROM golang:1.17
WORKDIR /mnt
    # install deps
    # RUN go get -d -v a/go/package/name
# copy source
COPY ./src .
# build
RUN mkdir -p dist
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dist/app .

# --- Run Stage

FROM alpine:latest  
RUN apk --no-cache add ca-certificates

# copy binary from build stage
WORKDIR /root/
COPY --from=0 /mnt/dist/ ./

# when the container is started, run the binary
CMD ["./app"]

