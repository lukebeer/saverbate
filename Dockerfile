FROM golang:1.14.3-buster AS compile

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY ./ ./

ARG SERVICE_NAME

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/service cmd/${SERVICE_NAME}/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add tini curl

ENTRYPOINT ["tini", "--"]

COPY --from=compile /app/service /app/
