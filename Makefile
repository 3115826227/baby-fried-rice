GO=CGO_ENABLED=0 GO111MODULE=on go
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES= \
	cmd/gateway/gateway \
	cmd/userAccount/userAccount \
	cmd/adminAccount/adminAccount \
	cmd/rootAccount/rootAccount \
	cmd/accountDao/accountDao \
	cmd/spaceDao/spaceDao \
	cmd/space/space

build: $(MICROSERVICES)

cmd/gateway/gateway:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/gateway

cmd/userAccount/userAccount:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/userAccount

cmd/adminAccount/adminAccount:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/adminAccount

cmd/rootAccount/rootAccount:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/rootAccount

cmd/accountDao/accountDao:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/accountDao

cmd/spaceDao/spaceDao:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/spaceDao

cmd/space/space:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/space

clean:
	rm -f $(MICROSERVICES)

run:
	cd bin && ./baby-launch.sh