default : fmt test book

test :
	go test ./...

fmt :
	find . -name \*.go | xargs gofmt -w

book: doc/book/bachdoc/null.md doc/book/bachdoc/io.md doc/book/bachdoc/logic.md doc/book/bachdoc/math.md doc/book/bachdoc/text.md doc/book/bachdoc/arr.md doc/book/bachdoc/obj.md doc/book/bachdoc/types.md doc/book/bachdoc/values.md doc/book/bachdoc/regexp.md doc/book/bachdoc/control.md
	mdbook build doc/book

doc/book/bachdoc/%.md : builtin/%.go
	mkdir -p $$(dirname $@)
	go run bachdoc/main.go $* > $@
