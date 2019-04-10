fmt :
	for file in *.go */*.go; do go fmt $$file; done
