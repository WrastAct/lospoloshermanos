FROM golang:1.19.2

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
ENV PORT=4000\
    MAESTRO_DB_DSN=ASA

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]

