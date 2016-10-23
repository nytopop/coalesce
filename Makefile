default: build

build:
	go fmt *.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o bin/coalesce

docker:
	go fmt *.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o bin/coalesce
	sudo docker-compose build

run:
	sudo docker-compose up
