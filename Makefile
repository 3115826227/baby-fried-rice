GO=CGO_ENABLED=0 GO111MODULE=on go
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES= \
	cmd/gateway/gateway \
	cmd/userAccount/userAccount \
	cmd/rootAccount/rootAccount \
	cmd/accountDao/accountDao \
	cmd/spaceDao/spaceDao \
	cmd/space/space \
	cmd/imDao/imDao \
	cmd/im/im \
	cmd/connect/connect \
	cmd/file/file

build: $(MICROSERVICES)

cmd/gateway/gateway:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/gateway

cmd/userAccount/userAccount:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/userAccount

cmd/rootAccount/rootAccount:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/rootAccount

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


clean:
	rm -f $(MICROSERVICES)

run:
	cd bin && ./baby-launch.sh