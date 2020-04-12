APP = trackers
BUILD_DIR ?= build

default: build

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: $(BUILD_DIR)/$(APP)
$(BUILD_DIR)/$(APP):
	go build -mod=vendor -v -o $(BUILD_DIR)/$(APP) cmd/$(APP)/*

build: $(BUILD_DIR)/$(APP)

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR) > /dev/null

.PHONY: clean-vendor
clean-vendor:
	rm -rf ./vendor > /dev/null

.PHONY: test
test:
	go test -cover -v pkg/$(APP)/*

install: clean build
	cp -v -f $(BUILD_DIR)/$(APP) /usr/local/bin/$(APP)
	chmod +x /usr/local/bin/$(APP)

tag:
	git tag `cat version`
