GO=CGO_ENABLED=0 GO111MODULE=on go
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES= \
	cmd/gateway/gateway \
	cmd/manage/manage \
	cmd/userAccount/userAccount \
	cmd/accountDao/accountDao \
	cmd/spaceDao/spaceDao \
	cmd/space/space \
	cmd/imDao/imDao \
	cmd/im/im \
	cmd/connect/connect \
	cmd/file/file \
	cmd/shopDao/shopDao \
	cmd/shop/shop \
	cmd/smsDao/smsDao \
	cmd/gameDao/gameDao \
    cmd/game/game \
    cmd/liveDao/liveDao \
    cmd/live/live \
    cmd/blogDao/blogDao \
    cmd/blog/blog

build: $(MICROSERVICES)

cmd/gateway/gateway:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/gateway

cmd/manage/manage:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/manage

cmd/userAccount/userAccount:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/userAccount

cmd/accountDao/accountDao:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/accountDao

cmd/spaceDao/spaceDao:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/spaceDao

cmd/space/space:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/space

cmd/imDao/imDao:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/imDao

cmd/im/im:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/im

cmd/connect/connect:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/connect

cmd/file/file:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/file

cmd/shopDao/shopDao:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/shopDao

cmd/shop/shop:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/shop

cmd/smsDao/smsDao:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/smsDao

cmd/gameDao/gameDao:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/gameDao

cmd/game/game:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/game

cmd/liveDao/liveDao:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/liveDao

cmd/live/live:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/live

cmd/blogDao/blogDao:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/blogDao

cmd/blog/blog:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/blog

clean:
	rm -f $(MICROSERVICES)

run:
	cd bin && ./baby-launch.sh