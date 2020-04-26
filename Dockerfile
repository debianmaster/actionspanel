FROM golang:1.14 as backend

# Install dependencies
WORKDIR /go/src/github.com/phunki/actionspanel
COPY ./go.mod ./go.sum ./
RUN go mod download

# Build artifacts
WORKDIR /go/src/github.com/phunki/actionspanel
COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY ./mock ./mock
RUN go get github.com/golang/mock/mockgen
RUN go generate ./...
RUN go test -race ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/actionspanel cmd/actionspanel.go

FROM scratch
# SSL Certs
COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=backend /go/bin/actionspanel /go/bin/actionspanel

# Run the hello binary.
ENTRYPOINT ["/go/bin/actionspanel", "api"]
