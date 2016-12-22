PACKAGE = github.com/cavaliercoder/redirector

redirector_SOURCES = \
	config.go \
	database.go \
	gob.go \
	keybuilder.go \
	logger.go \
	main.go \
	management.go \
	management_client.go \
	mapping.go \
	runtime.go \
	server.go

all: redirector

redirector: $(redirector_SOURCES)
	go build -x -o redirector $(redirector_SOURCES)

clean:
	go clean -x -i $(PACKAGE)
