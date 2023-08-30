FROM golang:alpine AS builder

WORKDIR /app
COPY . .

RUN apk update && apk add --no-cache git make
RUN go mod tidy
RUN make server

FROM scratch

# copy over ssl ca certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder app/bin/server bin/server
EXPOSE 8081

# Run the binary
CMD ["bin/server", "--port=8081"]
