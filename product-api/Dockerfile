FROM golang:1.21-alpine

WORKDIR /app

COPY . .

ENV GO111MODULE=on

RUN go mod download

RUN go build -o main .

CMD ["./main"]