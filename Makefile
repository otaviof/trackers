APP = trackers
BUILD_DIR ?= build

.PHONY: default bootstrap build clean test

default: build

bootstrap:
	dep ensure -v -vendor-only

build:
	go build -v -o $(BUILD_DIR)/$(APP) cmd/$(APP)/*

clean:
	rm -rf $(BUILD_DIR) > /dev/null

clean-vendor:
	rm -rf ./vendor > /dev/null

test:
	go test -cover -v pkg/$(APP)/*

install: clean build
	cp -v -f $(BUILD_DIR)/$(APP) /usr/local/bin/$(APP)
	chmod +x /usr/local/bin/$(APP)

tag:
	git tag `cat version`
