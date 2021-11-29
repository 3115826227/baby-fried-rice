FROM golang:1.15 AS build
WORKDIR /baby-fried-rice
COPY . .
RUN go env -w GOPROXY=https://goproxy.io && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o blog_dao cmd/blogDao/blog_dao.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/blogDao
RUN mkdir res && mkdir cert && mkdir -p /baby/logs/
COPY --from=build /baby-fried-rice/cmd/blogDao/res/config.yaml res/
COPY --from=build /baby-fried-rice/cmd/blogDao/cert cert/
COPY --from=build /baby-fried-rice/blogDao .

EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./blog_dao"]