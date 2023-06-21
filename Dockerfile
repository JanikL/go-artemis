FROM golang:1.20-alpine
WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
COPY artemis ./artemis
COPY examples ./examples
RUN go build -o artemis_example examples/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=0 /app/artemis_example ./

CMD ["./artemis_example"]