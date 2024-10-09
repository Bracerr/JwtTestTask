FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o my-app ./src/cmd/main.go

FROM alpine:3.17

RUN apk update && apk --no-cache add ca-certificates

COPY --from=builder /app/my-app /usr/local/bin/my-app

RUN chmod +x /usr/local/bin/my-app

#COPY .env ./

EXPOSE 8080

CMD ["my-app"]
