.PHONY: test cover

test:
	go test -fullpath ./...

cover:
	$(eval COV_FILE := $(shell mktemp))
	go test -v ./... -coverprofile $(COV_FILE)
	$(eval COV_HTML_FILE := $(shell mktemp))
	go tool cover -html $(COV_FILE) -o $(COV_HTML_FILE)
	@rm $(COV_FILE)
	@firefox $(COV_HTML_FILE)
	@sleep 0.5 # wait for firefox to load file
	@rm $(COV_HTML_FILE)
