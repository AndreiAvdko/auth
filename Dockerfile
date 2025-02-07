FROM golang:1.20.3-alpine AS builder

COPY . /github.com/AndreiAvdko/auth/source/
WORKDIR /github.com/AndreiAvdko/auth/source/


RUN mod download
RUN go build -o ./bin/auth_server cmd/main.go


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/AndreiAvdko/auth/source/bin/auth_server .

CMD ["./auth_server"]