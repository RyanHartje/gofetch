FROM golang:1.9-alpine

RUN mkdir -p /opt/gold
COPY /opt/gold/

WORKDIR /opt/gold
RUN go get

CMD go run /opt/gold/server/main.go /opt/gold/server/auth.go
