FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

EXPOSE 8081

COPY . .

RUN export GO111MODULE=on
RUN go mod tidy
RUN go build -o binary server/main.go

ENTRYPOINT ["/app/binary"]