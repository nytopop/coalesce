# We statically link because cannot link to scratch container's libc.

# Build .deb package.
bin/coalesce.deb: bin/coalesce
	echo building deb package!

# Build binary
bin/coalesce:
	go fmt *.go
	go build -ldflags '-linkmode external -extldflags -static -w' -o bin/coalesce

clean:
	rm -f bin/coalesce bin/coalesce.deb
