UNIT_TEST_FILES := $$(go list ./... | grep -v ./internal/server)

generate-repository-mocks: 
	cd internal && mockgen --source repository/main.go --destination mocks/repository.go -package mocks

# "-count=1" flag prevents caching of test results
run-server-e2e-tests: 
	GOFLAGS="-count=1" godotenv -f .env.test go test ./internal/server 

# run unit tests 
# exclude the server folder containing integration tests
run-unit-tests:
	@eval GOFLAGS="-count=1" go test $(UNIT_TEST_FILES) -cover

# generate swagger specs 
swagger-gen:
	go mod vendor \
		&& rm -rf internal/docs \
		&& swag init -d internal/server,internal/handlers,internal/dto,internal/models,internal/utils -o internal/docs --generatedTime=true

# build docker image for server 
build-server-docker-image: 
	docker build -f Dockerfile.server -t keeper-go . 

# build docker image for client 
build-client-docker-image: 
	docker build -f Dockerfile.client -t keeper-client . 

build-docker-compose: 
	docker-compose up --build 

run-docker-compose: 
	docker-compose up 

stop-docker-compose: 
	docker-compose down

start-server:
	go run cmd/server.go

start-client-dev:
	cd client && npm run dev

.PHONY: run-unit-tests