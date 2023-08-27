.PHONY: docker
docker:
	@rm ebook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o ebook .
	@docker rmi -f stone/ebook-live:v0.0.1
	@docker build -t stone/ebook-live:v0.0.1 .