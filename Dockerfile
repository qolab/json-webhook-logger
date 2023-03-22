FROM golang:1.18-alpine AS build

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -v -ldflags="-w -s" -o ./webhook-json-logger

FROM alpine:3.17

ENV USER_ID=65535
ENV GROUP_ID=65535
ENV USER_NAME=webhook
ENV GROUP_NAME=webhook

RUN addgroup -g $GROUP_ID $GROUP_NAME && \
    adduser --shell /sbin/nologin --disabled-password \
    -h /app --uid $USER_ID --ingroup $GROUP_NAME $USER_NAME

RUN apk add --no-cache ca-certificates

COPY --from=build /app/webhook-json-logger /app/webhook-json-logger

USER webhook
WORKDIR /app

EXPOSE 8080

CMD ["/app/webhook-json-logger"]
