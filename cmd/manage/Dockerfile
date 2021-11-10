FROM golang:1.15 AS build
WORKDIR /go/src/baby-fried-rice
COPY . .
RUN go env -w GO111MODULE=off
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o manage cmd/manage/manage.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/manage
RUN mkdir res && mkdir -p /baby/logs/
COPY --from=build /go/src/baby-fried-rice/cmd/manage/res/config.yaml res/
COPY --from=build /go/src/baby-fried-rice/manage .

EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./manage"]