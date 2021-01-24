OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
CUR := $(shell pwd)
IMAGE_NAME := helm-diff-notifier
IMAGE := hanks/helm-diff-notifier:dev
IMAGE_VERSION := $(shell grep -o '[[:digit:]].[[:digit:]].[[:digit:]]' version/version.go)
WORKSPACE := /go/src/github.com/hanks/helm-diff-notifier

test: lint-go
	docker run -it --rm -v $(CUR):${WORKSPACE} ${IMAGE} \
		go test -v -covermode=count -coverprofile=coverage.out ./...
run:
	docker run -it --rm -v $(CUR):${WORKSPACE} ${IMAGE} \
		go run main.go

lint-go:
	docker run -it --rm -v $(CUR):${WORKSPACE} ${IMAGE} \
		staticcheck -f json ./...

lint-dockerfile:
	docker run -it --rm -v $(CUR):${WORKSPACE} ${IMAGE} \
		hadolint Dockerfile

build-dev-image: lint-dockerfile
	docker build -t ${IMAGE} .

push: build-dev-image
	docker push ${IMAGE}

build: clean
	docker run -it --rm -v $(CUR):${WORKSPACE} -e "CGO_ENABLED=0" -e "GOARCH=amd64" -e "GOOS=linux" ${IMAGE} go build -o ./dist/bin/${IMAGE_NAME}_linux_amd64_${IMAGE_VERSION} main.go
	sudo zip -j ./dist/bin/$(IMAGE_NAME)_linux_amd64_$(IMAGE_VERSION).zip ./dist/bin/$(IMAGE_NAME)_linux_amd64_$(IMAGE_VERSION)
	sudo rm ./dist/bin/$(IMAGE_NAME)_linux_amd64_$(IMAGE_VERSION)

	docker run -it --rm -v $(CUR):${WORKSPACE} -e "CGO_ENABLED=0" -e "GOARCH=amd64" -e "GOOS=darwin" ${IMAGE} go build -o ./dist/bin/${IMAGE_NAME}_darwin_amd64_${IMAGE_VERSION} main.go
	sudo zip -j ./dist/bin/$(IMAGE_NAME)_darwin_amd64_$(IMAGE_VERSION).zip ./dist/bin/$(IMAGE_NAME)_darwin_amd64_$(IMAGE_VERSION)
	sudo rm ./dist/bin/$(IMAGE_NAME)_darwin_amd64_$(IMAGE_VERSION)

coveralls:
	docker run -it --rm -v $(CUR):${WORKSPACE} ${IMAGE} goveralls -coverprofile=coverage.out -service=circle-ci -repotoken ${COVERALLS_TOKEN}

clean:
	rm -rf ./dist

install:
	unzip -o ./dist/bin/$(IMAGE_NAME)_$(OS)_amd64_$(IMAGE_VERSION).zip -d ./dist/bin/
	cp ./dist/bin/$(IMAGE_NAME)_$(OS)_amd64_$(IMAGE_VERSION) /usr/local/bin/$(IMAGE_NAME)
	rm ./dist/bin/$(IMAGE_NAME)_$(OS)_amd64_$(IMAGE_VERSION)

uninstall:
	rm /usr/local/bin/$(IMAGE_NAME)

changelog:
	git-chglog -o CHANGELOG.md
