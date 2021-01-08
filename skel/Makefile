APPS     :=	sonic

all:	$(APPS)

$(APPS):
	go build -o bin/$@ cmd/$@/*.go

asset:
	cd client && \
		yarn install && \
		./scripts/build.sh

clean:
	rm -rf bin/*
