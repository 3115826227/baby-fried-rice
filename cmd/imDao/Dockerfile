FROM golang:1.15 AS build
WORKDIR /baby-fried-rice
COPY . .
RUN go env -w GOPROXY=https://goproxy.io
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o im_dao cmd/imDao/im_dao.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/imDao
RUN mkdir res && mkdir -p /baby/logs/
COPY --from=build /baby-fried-rice/cmd/imDao/res/config.yaml res/
COPY --from=build /baby-fried-rice/cmd/accountDao/cert cert/
COPY --from=build /baby-fried-rice/im_dao .
EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./im_dao"]
