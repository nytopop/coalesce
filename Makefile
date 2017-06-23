PNAME = coalesce_0.1.1-1
BDIR = build
PDIR = $(BDIR)/$(PNAME)

# Build .deb packages.
debian.amd64: test amd64
	# Binary
	mkdir -p $(PDIR)/usr/bin
	cp $(BDIR)/coalesce-amd64 $(PDIR)/usr/bin/coalesce
	# Config
	mkdir -p $(PDIR)/etc
	cp defaults/coalesce.conf $(PDIR)/etc/coalesce.conf
	# Logging
	mkdir -p $(PDIR)/var/log/coalesce
	# Database
	mkdir -p $(PDIR)/var/lib/coalesce
	cp defaults/init.sql $(PDIR)/var/lib/coalesce/init.sql
	# Resources
	mkdir -p $(PDIR)/usr/share/coalesce
	cp -r resources/* $(PDIR)/usr/share/coalesce/
	# systemd
	mkdir -p $(PDIR)/etc/systemd/system
	cp defaults/coalesce.service $(PDIR)/etc/systemd/system/coalesce.service
	# Packaging
	mkdir $(PDIR)/DEBIAN
	cp debian/* $(PDIR)/DEBIAN/
	dpkg-deb --build $(PDIR)

# Run tests
test:
	go test ./...

# We statically link for scratch container support.
amd64: clean
	go fmt *.go
	env GOOS=linux GOARCH=amd64 go build -ldflags '-linkmode external -extldflags -static -w' -o $(BDIR)/coalesce-amd64

clean:
	rm -rf build

run: amd64
	rm -rf run
	mkdir run
	build/coalesce-amd64 -cfg defaults/testing.conf

prun: amd64
	build/coalesce-amd64 -cfg defaults/testing.conf
