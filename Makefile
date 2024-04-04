.PHONY: clean
clean:
	@echo "==> Cleaning releases"
	@GOOS=linux go clean -i -x ./...

.PHONY: build
build: compile
	@echo "==> Building docker image digitalocean/kubelet-rubber-stamp:$(VERSION)"
	@docker build -t digitalocean/kubelet-rubber-stamp:$(VERSION)" .

.PHONY: push
push:
	@echo "==> Publishing digitalocean/kubelet-rubber-stamp:$(VERSION)"
	@docker push digitalocean/kubelet-rubber-stamp:$(VERSION)
	@echo "==> Your image is now available at digitalocean/kubelet-rubber-stamp:$(VERSION)"

publish: clean build push