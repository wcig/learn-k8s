FROM golang:1.24.5-alpine3.22 as builder
WORKDIR /goapp
#COPY go.mod go.sum main.go controller.go ./
COPY go.mod go.sum ./
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go mod download
COPY . .

RUN ls -la

RUN CGO_ENABLED=0 GOOS=linux go build -o ingress-manager-example .

FROM alpine:3.22
WORKDIR /goapp
COPY --from=builder /goapp/ingress-manager-example .
EXPOSE 8080
CMD ["./ingress-manager-example"]