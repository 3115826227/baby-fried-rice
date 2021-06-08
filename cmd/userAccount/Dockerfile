FROM golang:1.15 AS build
WORKDIR /baby-fried-rice
COPY . .
RUN go env -w GOPROXY=https://goproxy.io
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags=jsoniter -o user_account cmd/userAccount/user_account.go

FROM alpine:3.8
ENV TZ=Asia/Shanghai LANG=C.UTF-8 TERM=xterm
WORKDIR /go/src/baby-fried-rice/module/userAccount
RUN mkdir res && mkdir cert && mkdir -p /baby/logs/
COPY --from=build /baby-fried-rice/cmd/userAccount/res/config.yaml res/
COPY --from=build /baby-fried-rice/cmd/userAccount/cert cert/
COPY --from=build /baby-fried-rice/user_account .
EXPOSE 8080
VOLUME ["/baby/logs"]
ENTRYPOINT ["./user_account"]