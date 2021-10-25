FROM golang:1.15 AS build
WORKDIR /go/src/baby-fried-rice
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o space_dao cmd/spaceDao/space_dao.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/spaceDao
RUN mkdir res && mkdir -p /baby/logs/
COPY --from=build /go/src/baby-fried-rice/cmd/spaceDao/res/config.yaml res/
COPY --from=build /go/src/baby-fried-rice/cmd/spaceDao/cert cert/
COPY --from=build /go/src/baby-fried-rice/space_dao .
EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./space_dao"]
