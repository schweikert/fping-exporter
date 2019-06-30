VERSION = 0.1

build:
	gox -os="linux darwin windows" -arch="amd64" -verbose \
	    -ldflags "-X main.buildCommit=`git rev-parse --short HEAD` \
	              -X main.buildDate=`date +%Y-%m-%d` \
	              -X main.buildVersion=$(VERSION)" \
	    ./...

before_build:
	go get github.com/mitchellh/gox
	
lint:
	golangci-lint run *.go
