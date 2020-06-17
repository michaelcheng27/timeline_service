build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/timeline src/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/s3_handler src/s3_handler/main.go
	
.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
