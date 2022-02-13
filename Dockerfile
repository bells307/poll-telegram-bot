FROM golang:latest

WORKDIR /usr/src/poll-telegram-bot

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/poll-telegram-bot ./cmd

CMD ["poll-telegram-bot"]