PACKAGE = github.com/cavaliercoder/redirector

redirector_SOURCES = \
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
	runtime.go

all: redirector

redirector: $(redirector_SOURCES)
	go build -x -o redirector $(redirector_SOURCES)

clean:
	go clean -x -i $(PACKAGE)

.PHONY: all clean
