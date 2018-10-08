# Stage 1. Build the binary
FROM golang:1.11

# add a non-privileged user
RUN useradd -u 10001 myapp

RUN mkdir -p /go/src/github.com/nikk-dzhurov/go-workshop
ADD . /go/src/github.com/nikk-dzhurov/go-workshop
WORKDIR /go/src/github.com/nikk-dzhurov/go-workshop

# build the binary with go build
RUN CGO_ENABLED=0 go build \
	-o bin/go-workshop github.com/nikk-dzhurov/go-workshop/cmd/go-workshop

# Stage 2. Run the binary
FROM scratch

ENV PORT 8080
ENV DIAG_PORT 8585

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=0 /etc/passwd /etc/passwd
USER myapp

COPY --from=0 /go/src/github.com/nikk-dzhurov/go-workshop/bin/go-workshop /go-workshop
EXPOSE $PORT
EXPOSE $DIAG_PORT

CMD ["/go-workshop"]
