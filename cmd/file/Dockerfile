FROM golang:1.15 AS build
WORKDIR /go/src/baby-fried-rice
COPY . .
RUN go env -w GO111MODULE=off
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o file cmd/file/file.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/file
RUN mkdir res && mkdir cert && mkdir -p /baby/logs/
COPY --from=build /go/src/baby-fried-rice/cmd/file/res/config.yaml res/
COPY --from=build /go/src/baby-fried-rice/file .

EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./file"]