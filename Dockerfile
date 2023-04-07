FROM golang:1.19-alpine

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o app examples/main.go

FROM alpine:latest
COPY --from=0 app .

CMD ["./app"]
