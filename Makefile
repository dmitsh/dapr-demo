#SRCS := $(wildcard cmd/*/*.go)
#PRGS := $(notdir $(patsubst %.go,%,$(SRCS)))
#BINS := $(addprefix bin/,$(PRGS))
DIRS := cmd/pubsub/pub cmd/pubsub/sub cmd/pubsub/pubsub cmd/statestore

.PHONY: build build-linux go-mod go-check clean docker-build docker-push

build:
	mkdir -p bin
	cd bin; for f in $(DIRS); do echo $$f; CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' ../$$f; done

build-linux:
	mkdir -p bin/linux
	cd bin/linux; for f in $(DIRS); do echo $$f; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' ../../$$f; done

go-mod: go-check
	go mod tidy
	go mod vendor

go-check:
	@which go > /dev/null

clean:
	rm -rf ./bin

DOCKER_REPO=docker.io/dmitsh
DOCKER_IMAGE_VER=0.1
DOCKER_IMAGE=dapr-demo:${DOCKER_IMAGE_VER}

docker-build: build-linux
	docker build -t ${DOCKER_IMAGE} .

docker-push:
	docker tag ${DOCKER_IMAGE} ${DOCKER_REPO}/${DOCKER_IMAGE} && docker push ${DOCKER_REPO}/${DOCKER_IMAGE}
