APPS     :=	sonic

all:	$(APPS)

$(APPS):
	go build --ldflags "-X main.Version=$$(cat .tags)" -o bin/$@ app.go

test:
	go test ./...

asset:
	cd client && \
		yarn install && \
		./scripts/build.sh

clean:
	rm -rf bin/*
