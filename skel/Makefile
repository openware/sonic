APPS     :=	sonic

all:	$(APPS)

$(APPS):
	go build -o bin/$@ app.go

test:
	go test ./...

asset:
	cd client && \
		yarn install && \
		./scripts/build.sh

clean:
	rm -rf bin/*
