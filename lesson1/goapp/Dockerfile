FROM golang:1.24.5-alpine3.22 as builder
WORKDIR /goapp
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct
ADD main.go /goapp/main.go
ADD go.mod /goapp/go.mod
ADD go.sum /goapp/go.sum
RUN go build

FROM alpine:3.22
WORKDIR /goapp
RUN mkdir -p /goapp/log
COPY --from=builder /goapp/goapp .
EXPOSE 8080
CMD ["./goapp"]