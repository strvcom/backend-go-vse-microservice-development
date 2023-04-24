.PHONY: test test-fuzzy

test:
	go test -v ./...

test-fuzzy:
	 go test -v -fuzztime=5s -fuzz=.
