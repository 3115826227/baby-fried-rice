FROM golang:1.15 AS build
WORKDIR /baby-fried-rice
COPY . .
RUN go env -w GOPROXY=https://goproxy.io && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o game_dao cmd/gameDao/game_dao.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/gameDao
RUN mkdir res && mkdir cert && mkdir -p /baby/logs/
COPY --from=build /baby-fried-rice/cmd/gameDao/res/config.yaml res/
COPY --from=build /baby-fried-rice/cmd/gameDao/cert cert/
COPY --from=build /baby-fried-rice/gameDao .

EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./game_dao"]