default: build

build:
	go fmt *.go
	go build -ldflags '-linkmode external -extldflags -static -w' -o bin/coalesce

docker:
	sudo docker build -t nytopop/coalesce .
#	sudo docker-compose build

run:
	sudo docker run nytopop/coalesce
#	sudo docker-compose up

all:
	go fmt *.go
	go build -ldflags '-linkmode external -extldflags -static -w' -o bin/coalesce
	sudo docker build -t nytopop/coalesce .
	sudo docker run nytopop/coalesce
#	sudo docker-compose build
#	sudo docker-compose up
