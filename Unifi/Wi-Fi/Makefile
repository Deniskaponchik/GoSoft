.PHONY: build
build:
	go build -o bin\$(CNTRL)_v$(VERSION) -ldflags="-s -w"
    start /B ./bin/$(CNTRL)_v$(VERSION)

.PHONY: test
test:
	go build -o bin\$(CNTRL)_v$(VERSION)-TEST -ldflags="-s -w"
	start /B ./bin/$(CNTRL)_v$(VERSION)-TEST

.DEFAULT_GOAL := test