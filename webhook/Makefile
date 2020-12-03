IMAGE_NAME := "rpkatz/nullsolver-webhook"
IMAGE_TAG := "v1.1.1"

build:
	CGO_ENABLED=0 go build -o output/webhook -ldflags '-w -extldflags "-static"' .
	docker build -t "$(IMAGE_NAME):$(IMAGE_TAG)" .
