.PHONY: docker
docker:
	@rm ebook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o ebook .
	@docker rmi -f stone/ebook-live:v0.0.1
	@docker build -t stone/ebook-live:v0.0.1 .


.PHONY: mock
mock:
	@mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/service/user.go -package=svcmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/service/mocks/user.mock.go
	@mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/service/code.go -package=svcmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/service/mocks/code.mock.go

	@mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/user.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/mocks/user.mock.go
	@mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/code.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/mocks/code.mock.go

	@mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/dao/user.go -package=daomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/dao/mocks/user.mock.go

	@mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/cache/user.go -package=cachemocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/cache/code.go -package=cachemocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/cache/mocks/code.mock.go

	@mockgen -source=ebook/pkg/ratelimit/types.go -package=limitmocks -destination=webook/pkg/ratelimit/mocks/ratelimit.mock.go
	@mockgen -package=redismocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/cache/cachemocks/cmdable.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy