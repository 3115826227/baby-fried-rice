FROM golang:1.15 AS build
WORKDIR /baby-fried-rice
COPY . .
RUN go env -w GOPROXY=https://goproxy.io && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o shop_dao cmd/shopDao/shop_dao.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/shopDao
RUN mkdir res && mkdir -p /baby/logs/
COPY --from=build /baby-fried-rice/cmd/shopDao/res/config.yaml res/
COPY --from=build /baby-fried-rice/cmd/shopDao/cert cert/
COPY --from=build /baby-fried-rice/shop_dao .
EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./shop_dao"]
