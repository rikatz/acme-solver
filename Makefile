IMAGE_NAME := "rpkatz/acme-solver"
IMAGE_TAG := "v1.1.1"

build:
	CGO_ENABLED=0 go build -o output/acme-solver -ldflags '-w -extldflags "-static"' .
	docker build -t "$(IMAGE_NAME):$(IMAGE_TAG)" .
