GO = go
GOFILES = wiki.go
BINARY = wiki

DOCKER_IMAGE = go-wiki
DOCKER_CONTAINER = go-wiki-container
DOCKER_PORT = 8080

# targets
.PHONY: all run build clean docker-build docker-run docker-clean docker-stop

# default target - run the Go app locally
all: run

# build the Go application locally
build:
	$(GO) build -o $(BINARY) $(GOFILES)

# run the Go application locally
run:
	$(GO) run $(GOFILES)

# clean up Go build files
clean:
	$(GO) clean
	rm -f $(BINARY)

# build the Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# run the Docker container
docker-run: docker-stop
	docker run -p $(DOCKER_PORT):8080 -v $(PWD)/data:/app/data --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)

# stop and remove any existing Docker container
docker-stop:
	@docker stop $(DOCKER_CONTAINER) 2>/dev/null || true
	@docker rm $(DOCKER_CONTAINER) 2>/dev/null || true

# clean up Docker images and containers
docker-clean: docker-stop
	docker rmi $(DOCKER_IMAGE)

# rebuild and run Docker container
docker-rebuild: docker-clean docker-build docker-run

