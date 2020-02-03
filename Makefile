.PHONY: test docker

DOCKER_IMG = cyrilix/robocar-pca9685

test:
	go test -race -mod vendor ./cmd/rc-pca9685 ./actuator ./part ./util

docker:
	docker buildx build . --platform linux/arm/7,linux/arm64,linux/amd64 -t ${DOCKER_IMG} --push

