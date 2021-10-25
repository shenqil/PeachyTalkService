##
## Build
##
FROM golang:1.16 AS build

WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY . .

RUN go mod download

RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -o /gin-admin

##
## Deploy
##

FROM fssq/distroless_base-debian10
# FROM scratch

WORKDIR /

# 复制编译后的程序
COPY --from=build /gin-admin /gin-admin

EXPOSE 8080

ENTRYPOINT ["/gin-admin"]