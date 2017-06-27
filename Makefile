PNAME = coalesce_0.1.1-1
BDIR = build
PDIR = $(BDIR)/$(PNAME)

default: amd64

test:
	go test ./...

amd64: clean
	go fmt *.go
	env GOOS=linux GOARCH=amd64 go build -o $(BDIR)/coalesce-amd64

clean:
	rm -rf build

run: amd64
	rm -rf run
	mkdir run
	build/coalesce-amd64 -cfg defaults/testing.conf

prun: amd64
	build/coalesce-amd64 -cfg defaults/testing.conf
