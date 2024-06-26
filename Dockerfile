FROM golang:1.21-alpine AS builder

# Set the current working directory.
RUN mkdir -p /avataro
WORKDIR /avataro

# install the modules
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

# Build executable binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o avataro-app app/main.go

# Remove unnecessary folders
RUN rm -Rf tools

# Build the final image
FROM alpine:3.20

# Set the current working directory.
WORKDIR /avataro
COPY --from=builder /avataro /avataro

EXPOSE 8080

CMD ["./avataro-app"]
