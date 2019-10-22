default : fmt test book

test :
	go test ./...

fmt :
	find . -name \*.go | xargs gofmt -w

book:
	mdbook build doc/book
