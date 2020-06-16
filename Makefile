build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/timeline src/main.go
	
.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
