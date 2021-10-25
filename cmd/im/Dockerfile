FROM golang:1.15 AS build
WORKDIR /go/src/baby-fried-rice
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o im cmd/im/im.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/im
RUN mkdir res && mkdir -p /baby/logs/
COPY --from=build /go/src/baby-fried-rice/cmd/im/res/config.yaml res/
COPY --from=build /go/src/baby-fried-rice/cmd/im/cert cert/
COPY --from=build /go/src/baby-fried-rice/im .
EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./im"]
