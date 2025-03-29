default : fmt tidy test book

fmt :
	find . -name \*.go | xargs gofmt -w

tidy :
	go mod tidy

test :
	go test ./...

book : doc/book/bachdoc/builtin/null.md doc/book/bachdoc/builtin/io.md doc/book/bachdoc/builtin/logic.md doc/book/bachdoc/builtin/math.md doc/book/bachdoc/builtin/text.md doc/book/bachdoc/builtin/arr.md doc/book/bachdoc/builtin/obj.md doc/book/bachdoc/builtin/types.md doc/book/bachdoc/builtin/values.md doc/book/bachdoc/builtin/regexp.md doc/book/bachdoc/builtin/control.md
	mdbook build doc/book

doc/book/bachdoc/builtin/%.md : builtin/%.go bachdoc/main.go
	mkdir -p "$$(dirname $@)"
	go run ./bachdoc builtin $* > $@

deploy : book
	if [ -z $$DPLDEST ]; then echo DPLDEST is unset; exit 1; fi
	rsync -Pahz doc/book/book/ $$DPLDEST/
