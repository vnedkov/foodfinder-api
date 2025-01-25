# foodfinder-api
Backend service for FoodFinder web site - a search engine for UMASS dining halls

## Building binary
```shell
go build .
```

## Building a Docker image
```shell
docker build --pull --rm -f Dockerfile -t foodfinderapi:latest .
```

## Unit Testing
```sh
go test ./...
```