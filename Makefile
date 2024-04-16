NAME ?= kubelet-rubber-stamp
VERSION ?= $(shell cat VERSION)
REGISTRY ?= digitalocean
DOCKER_IMAGE ?= $(REGISTRY)/$(NAME):$(VERSION)

GO_VERSION ?= $(shell go mod edit -print | grep -E '^go [[:digit:].]*' | cut -d' ' -f2)

.PHONY: clean
clean:
	@echo "==> Cleaning releases"
	@GOOS=linux go clean -i -x ./...

.PHONY: build
build: compile
	@echo "==> Building docker image $(DOCKER_IMAGE)"
	@docker build --build-arg GOVERSION=$(GO_VERSION) -t $(DOCKER_IMAGE) .

.PHONY: push
push:
	@echo "==> Publishing $(DOCKER_IMAGE)"
	@docker push $(DOCKER_IMAGE)
	@echo "==> Your image is now available at $(DOCKER_IMAGE)"

.PHONY: publish
publish: clean build push

.PHONY: test
	@dev/test

.PHONY: bump-version
bump-version:
	@[ "${NEW_VERSION}" ] || ( echo "NEW_VERSION must be set (ex. make NEW_VERSION=vX.Y.Z bump-version)"; exit 1 )
	@(echo ${NEW_VERSION} | grep -E "^v") || ( echo "NEW_VERSION must be a semver ('v' prefix is required)"; exit 1 )
	@echo "Bumping VERSION from $(VERSION) to $(NEW_VERSION)"
	@echo $(NEW_VERSION) > VERSION