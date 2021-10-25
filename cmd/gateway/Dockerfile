FROM golang:1.15 AS build
WORKDIR /go/src/baby-fried-rice
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o gateway cmd/gateway/gateway.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/gateway
RUN mkdir res && mkdir -p /baby/logs/
COPY --from=build /go/src/baby-fried-rice/cmd/gateway/res/config.yaml res/
COPY --from=build /go/src/baby-fried-rice/gateway .
EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./gateway"]
