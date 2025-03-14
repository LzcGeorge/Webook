.PHONY: mock
mock:
	@mockgen -source=./webook/internal/service/user.go -package=svcmocks -destination=./webook/internal/service/mocks/user.mock.go
	@mockgen -source=./webook/internal/service/code.go -package=svcmocks -destination=./webook/internal/service/mocks/code.mock.go
	@mockgen -source=./webook/internal/service/article.go -package=svcmocks -destination=./webook/internal/service/mocks/article.mock.go
	@mockgen -source=./webook/internal/service/interactive.go -package=svcmocks -destination=./webook/internal/service/mocks/interactive.mock.go

	@mockgen -source=./webook/internal/repository/user.go -package=repomocks -destination=./webook/internal/repository/mocks/user.mock.go
	@mockgen -source=./webook/internal/repository/code.go -package=repomocks -destination=./webook/internal/repository/mocks/code.mock.go
	@mockgen -source=./webook/internal/repository/dao/user.go -package=daomocks -destination=./webook/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./webook/internal/repository/cache/user.go -package=cachemocks -destination=./webook/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=./webook/internal/repository/article/article.go -package=repomocks -destination=./webook/internal/repository/article/mocks/article.mock.go
	@mockgen -source=./webook/internal/repository/article/article_author.go -package=repomocks -destination=./webook/internal/repository/article/mocks/article_author.mock.go
	@mockgen -source=./webook/internal/repository/article/article_reader.go -package=repomocks -destination=./webook/internal/repository/article/mocks/article_reader.mock.go

	@mockgen -package=redismocks -destination=./webook/internal/repository/cache/redismocks/cmd.mock.go github.com/redis/go-redis/v9 Cmdable

	@go mod tidy
