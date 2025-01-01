.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f techselfknow/webook:v0.0.1 
	@docker build -t techselfknow/webook:v0.0.1 .
