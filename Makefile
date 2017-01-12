PACKAGE = redirector
PACKAGE_VERSION = 1.1.0
PACKAGE_PATH = github.com/cavaliercoder/$(PACKAGE)

SOURCES = \
	bolt.go \
	config.go \
	database.go \
	gob.go \
	handler.go \
	httperror.go \
	keybuilder.go \
	logger.go \
	main.go \
	management.go \
	management_client.go \
	mapping.go \
	redirect.go \
	redis.go \
	runtime.go \
	template.go

EXTRA_DIST = \
	Makefile \
	README.md

# see https://fedoraproject.org/wiki/PackagingDrafts/Go#Debuginfo
LDFLAGS = "-B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')"

all: $(PACKAGE)

$(PACKAGE): $(SOURCES)
	go build -a -v -x \
		-ldflags $(LDFLAGS) \
		-o $(PACKAGE) \
		$(SOURCES)

check:
	# for redis: $ docker run -d -p 6379:6379 redis
	go test -v -cover

clean:
	go clean -x -i .

install:
	go install $(PACKAGE_PATH)

get-deps:
	go get -v gopkg.in/urfave/cli.v1
	go get -v github.com/boltdb/bolt
	go get -v github.com/garyburd/redigo/redis

dist: $(SOURCES) $(EXTRA_DIST)
	rm -rvf $(PACKAGE)-$(PACKAGE_VERSION)/ || :
	mkdir $(PACKAGE)-$(PACKAGE_VERSION)/ || :
	cp -v $(SOURCES) $(EXTRA_DIST) $(PACKAGE)-$(PACKAGE_VERSION)/
	tar -zcvf $(PACKAGE)-$(PACKAGE_VERSION).tar.gz $(PACKAGE)-$(PACKAGE_VERSION)/
	rm -rvf $(PACKAGE)-$(PACKAGE_VERSION)/

distclean:
	rm -vf $(PACKAGE)-$(PACKAGE_VERSION).tar.gz

rpm: dist
	cp -v rpmbuild/$(PACKAGE).spec ~/rpmbuild/SPECS/$(PACKAGE).spec
	cp -v $(PACKAGE)-$(PACKAGE_VERSION).tar.gz ~/rpmbuild/SOURCES/$(PACKAGE)-$(PACKAGE_VERSION).tar.gz 
	rpmbuild -ba ~/rpmbuild/SPECS/$(PACKAGE).spec

rpmclean:
	rm -vf ~/rpmbuild/SPECS/$(PACKAGE).spec
	rm -vf ~/rpmbuild/SOURCES/$(PACKAGE)-$(PACKAGE_VERSION).tar.gz
	rm -vf ~/rpmbuild/RPMS/x86_64/redirector-*

.PHONY: all check clean install get-deps dist rpm
