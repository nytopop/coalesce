# We statically link because cannot link to scratch container's libc.

PNAME = coalesce_0.1-1
PDIR = build/$(PNAME)
BDIR = build

# Build .deb package.
bin/coalesce.deb: bin/coalesce
	mkdir $(PDIR)
	mkdir $(PDIR)/usr
	mkdir $(PDIR)/usr/bin
	cp $(BDIR)/coalesce $(PDIR)/usr/bin/coalesce
	mkdir $(PDIR)/DEBIAN
	cp debian/* $(PDIR)/DEBIAN/
	dpkg-deb --build $(PDIR)

# Build binary
bin/coalesce: clean
	go fmt *.go
	go build -ldflags '-linkmode external -extldflags -static -w' -o $(BDIR)/coalesce

clean:
	rm -rf build
