.PHONY: test-coverage test-cover-profile view-cover-profile save-cover-profile

test-coverage:
	go test -cover ./...

test-cover-profile:
	go test -coverprofile=coverage.out ./...

view-cover-profile:
	go tool cover -html=coverage.out

save-cover-profile:
	go tool cover -html=coverage.out -o test-coverage.html
