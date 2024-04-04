.PHONY: clean
clean:
	@echo "==> Cleaning releases"
	@GOOS=linux go clean -i -x ./...

.PHONY: build
build: compile
	@echo "==> Building docker image registry.digitalocean.com/ccp-infra/kubelet-rubber-stamp:$(VERSION)"
	@docker build -t registry.digitalocean.com/ccp-infra/kubelet-rubber-stamp:$(VERSION)" .

.PHONY: push
push:
	@echo "==> Publishing registry.digitalocean.com/ccp-infra/kubelet-rubber-stamp:$(VERSION)"
	@docker push registry.digitalocean.com/ccp-infra/kubelet-rubber-stamp:$(VERSION)
	@echo "==> Your image is now available at registry.digitalocean.com/ccp-infra/kubelet-rubber-stamp:$(VERSION)"

publish: clean build push