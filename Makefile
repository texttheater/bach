fmt :
	find . -name \*.go | parallel --gnu go fmt
