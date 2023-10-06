.PHONY: build
build:
	go build -o bin\Poly_v$(VERSION) -ldflags="-s -w"
    start /B ./bin/Poly_v$(VERSION)

.PHONY: test
test:
	go build -o bin\Poly_v$(VERSION)-TEST -ldflags="-s -w"
	start /B ./bin/Poly_v$(VERSION)-TEST

.DEFAULT_GOAL := test