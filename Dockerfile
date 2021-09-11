##
## Build
##
FROM golang:1.16-alpine

WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY go.sum .
COPY go.mod .
RUN go mod download

COPY *.go .

RUN go build -o /gin-admin

##
## Deploy
##

# FROM gcr.io/distroless/base-debian10
FROM scratch

WORKDIR /

COPY --from=build /gin-admin /gin-admin
COPY configs /configs

EXPOSE 8080

ENTRYPOINT ["/gin-admin"]