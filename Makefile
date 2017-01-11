default: build

build:
	go fmt *.go
	go build -ldflags '-linkmode external -extldflags -static -w' -o bin/coalesce

docker:
	sudo docker build -t nytopop/coalesce .

run:
	sudo docker run nytopop/coalesce

all:
	go fmt *.go
	go build -ldflags '-linkmode external -extldflags -static -w' -o bin/coalesce
	sudo docker build -t nytopop/coalesce .
	sudo docker run nytopop/coalesce
