APPS     :=	sonic

all:	$(APPS)

$(APPS):
	go build --ldflags "-X main.Version=$$(cat VERSION)" -o bin/$@ app.go

test:
	go test ./...

asset:
	cd client && \
		yarn install && \
		./scripts/build.sh && \
		cp -r build ../public/assets

clean:
	rm -rf bin/*


## line:14 we can't use build.sh in makefile or we need put in readme what build.sh should have