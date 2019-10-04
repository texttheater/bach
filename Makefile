default : fmt test

test :
	go test ./...

fmt :
	find . -name \*.go | xargs gofmt -w
